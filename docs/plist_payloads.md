# Payloads Flow: CREATE vs UPDATE

## **CREATE Flow - Payloads**

```
1. Get payloads string from Terraform schema
   └─> d.Get("payloads").(string)

2. HTML escape the entire string
   └─> html.EscapeString(payloads)

3. Assign directly to resource
   └─> resource.General.Payloads = escapedPayloads

4. Done - No further processing
```

**Simple and straightforward**: Just escape and assign.

---

## **UPDATE Flow - Payloads**

```
1. Fetch existing profile from Jamf Pro API
   └─> client.GetMacOSConfigurationProfileByID(resourceID)
   └─> Extract: existingProfile.General.Payloads (Jamf Pro's version with modified UUIDs)

2. Decode EXISTING payloads (from Jamf Pro)
   └─> plist.NewDecoder(existingPayload).Decode(&existingPlist)
   └─> Result: existingPlist map[string]any (contains Jamf Pro's UUIDs)

3. Decode NEW payloads (from Terraform state)
   └─> newPayload = d.Get("payloads").(string)
   └─> plist.NewDecoder(newPayload).Decode(&newPlist)
   └─> Result: newPlist map[string]any (contains user's original UUIDs)

4. Sync TOP-LEVEL UUIDs (Jamf Pro → Terraform)
   └─> newPlist["PayloadUUID"] = existingPlist["PayloadUUID"]
   └─> newPlist["PayloadIdentifier"] = existingPlist["PayloadIdentifier"]

5. Sync NESTED UUIDs
   └─> helpers.ExtractUUIDs(existingPlist, uuidMap, true)
   └─> helpers.ExtractPayloadIdentifiers(existingPlist, identifierMap, true)
   └─> helpers.UpdateUUIDs(newPlist, uuidMap, identifierMap, true)

6. Validate UUID matching
   └─> helpers.ValidatePayloadUUIDsMatch(existingPlist, newPlist, ...)
   └─> If mismatches found → return error

7. Re-encode the updated plist
   └─> encoder := plist.NewEncoder(&buf)
   └─> encoder.Indent("    ")
   └─> encoder.Encode(newPlist)
   └─> Result: buf contains XML plist with synced UUIDs

8. XML normalization (for embedding plist XML in Jamf Pro XML)
   └─> preMarshallingXMLPayloadUnescaping(buf.String())
       └─> Replace &#34; with "
   └─> preMarshallingXMLPayloadEscaping(unquotedContent)
       └─> Replace & with &amp;

9. Assign to resource
   └─> resource.General.Payloads = escapedContent
```

---

## **Why This Complexity?**

**Jamf Pro modifies UUIDs after creation**:
- When you CREATE a profile, Jamf Pro changes the top-level `PayloadUUID` and `PayloadIdentifier`
- Nested payload UUIDs stay the same
- If you UPDATE without preserving Jamf Pro's UUIDs, it treats it as a different profile
- This causes deployment problems on managed devices

**The UPDATE flow preserves UUID continuity** by:
1. Fetching what Jamf Pro actually stored
2. Extracting Jamf Pro's modified UUIDs
3. Injecting them into the new configuration
4. Re-encoding everything together

This ensures the update is recognized as a modification to the *same* profile, not a new profile.
