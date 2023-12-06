// computerinventory_data_source.go
package computerinventory

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProMacOSComputerInventory provides information about a specific computers inventory by its ID or Name.
func DataSourceJamfProComputerInventory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJamfProComputerInventoryRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"udid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"general": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lastIpAddress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lastReportedIp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"jamfBinaryVersion": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"platform": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"barcode1": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"barcode2": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"assetTag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remoteManagement": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"managed": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"managementUsername": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"supervised": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"mdmCapable": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"capable": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"capableUsers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"reportDate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lastContactTime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lastCloudBackupDate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lastEnrolledDate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mdmProfileExpiration": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"initialEntryDate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"distributionPoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enrollmentMethod": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"objectName": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"objectType": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"site": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"itunesStoreAccountActive": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enrolledViaAutomatedDeviceEnrollment": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"userApprovedMdm": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"declarativeDeviceManagementEnabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"extensionAttributes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definitionId": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"multiValue": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"dataType": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"options": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"inputType": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"managementId": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"diskEncryption": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bootPartitionEncryptionDetails": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"partitionName": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"partitionFileVault2State": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"partitionFileVault2Percent": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"individualRecoveryKeyValidityStatus": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"institutionalRecoveryKeyPresent": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"diskEncryptionConfigurationName": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fileVault2EnabledUserNames": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"fileVault2EligibilityMessage": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"purchasing": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"leased": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"purchased": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"poNumber": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"poDate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vendor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"warrantyDate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"appleCareId": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"leaseDate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"purchasePrice": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lifeExpectancy": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"purchasingAccount": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"purchasingContact": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extensionAttributes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definitionId": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"multiValue": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"dataType": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"options": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"inputType": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"macAppStore": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"sizeMegabytes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bundleId": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updateAvailable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"externalVersionId": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"storage": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bootDriveAvailableSpaceMegabytes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"device": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"model": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"revision": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"serialNumber": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sizeMegabytes": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"smartStatus": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"partitions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"sizeMegabytes": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"availableMegabytes": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"partitionType": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"percentUsed": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"fileVault2State": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"fileVault2ProgressPercent": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"lvmManaged": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"userAndLocation": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"realname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"position": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"phone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"departmentId": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"buildingId": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"room": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extensionAttributes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definitionId": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"multiValue": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"dataType": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"options": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"inputType": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"configurationProfiles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lastInstalled": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"removable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"displayName": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"profileIdentifier": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"printers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uri": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"hardware": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"make": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"model": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modelIdentifier": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"serialNumber": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"processorSpeedMhz": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"processorCount": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"coreCount": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"processorType": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"processorArchitecture": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"busSpeedMhz": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cacheSizeKilobytes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"networkAdapterType": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"macAddress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"altNetworkAdapterType": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"altMacAddress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"totalRamMegabytes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"openRamSlots": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"batteryCapacityPercent": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"smcVersion": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nicSpeed": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"opticalDrive": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bootRom": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bleCapable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"supportsIosAppInstalls": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"appleSilicon": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"extensionAttributes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definitionId": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"multiValue": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"dataType": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"options": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"inputType": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"localUserAccounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"userGuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fullName": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"admin": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"homeDirectory": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"homeDirectorySizeMb": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"fileVault2Enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"userAccountType": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"passwordMinLength": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"passwordMaxAge": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"passwordMinComplexCharacters": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"passwordHistoryDepth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"passwordRequireAlphanumeric": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"computerAzureActiveDirectoryId": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"userAzureActiveDirectoryId": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"azureActiveDirectoryId": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"commonName": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"identity": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"expirationDate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lifecycleStatus": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificateStatus": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subjectName": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"serialNumber": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sha1Fingerprint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"issuedDate": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fileType": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sizeBytes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"plugins": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"packageReceipts": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"installedByJamfPro": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"installedByInstallerSwu": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cached": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"fonts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"security": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sipStatus": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gatekeeperStatus": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"xprotectVersion": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"autoLoginDisabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"remoteDesktopEnabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"activationLockEnabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"recoveryLockEnabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"firewallEnabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"secureBootLevel": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"externalBootLevel": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bootstrapTokenAllowed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"operatingSystem": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"build": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"supplementalBuildVersion": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rapidSecurityResponse": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"activeDirectoryStatus": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fileVault2Status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"softwareUpdateDeviceId": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extensionAttributes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definitionId": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"multiValue": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"dataType": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"options": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"inputType": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"licensedSoftware": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ibeacons": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"softwareUpdates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"packageName": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"extensionAttributes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definitionId": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"multiValue": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"values": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"dataType": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"options": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"inputType": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"contentCaching": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"computerContentCachingInformationId": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parents": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"contentCachingParentId": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"alerts": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"contentCachingParentAlertId": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"addresses": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"className": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"postDate": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"details": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"contentCachingParentDetailsId": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"acPower": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"cacheSizeBytes": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"capabilities": {
													Type:     schema.TypeList,
													Computed: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"contentCachingParentCapabilitiesId": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"imports": {
																Type:     schema.TypeBool,
																Computed: true,
															},
															"namespaces": {
																Type:     schema.TypeBool,
																Computed: true,
															},
															"personalContent": {
																Type:     schema.TypeBool,
																Computed: true,
															},
															"queryParameters": {
																Type:     schema.TypeBool,
																Computed: true,
															},
															"sharedContent": {
																Type:     schema.TypeBool,
																Computed: true,
															},
															"prioritization": {
																Type:     schema.TypeBool,
																Computed: true,
															},
														},
													},
												},
												"portable": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"localNetwork": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"contentCachingParentLocalNetworkId": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"speed": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"wired": {
																Type:     schema.TypeBool,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"guid": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"healthy": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"alerts": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cacheBytesLimit": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"className": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"pathPreventingAccess": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"postDate": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"reservedVolumeBytes": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"resource": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"activated": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"active": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"actualCacheBytesUsed": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cacheDetails": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"computerContentCachingCacheDetailsId": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"categoryName": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"diskSpaceBytesUsed": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"cacheBytesFree": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cacheBytesLimit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cacheStatus": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cacheBytesUsed": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"dataMigrationCompleted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"dataMigrationProgressPercentage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"dataMigrationError": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"domain": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"userInfo": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeMap,
											Elem: &schema.Schema{
												Type: schema.TypeString,
											},
										},
									},
								},
							},
						},
						"maxCachePressureLast1HourPercentage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"personalCacheBytesFree": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"personalCacheBytesLimit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"personalCacheBytesUsed": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"publicAddress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"registrationError": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"registrationResponseCode": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"registrationStarted": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"registrationStatus": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"restrictedMedia": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"serverGuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"startupStatus": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tetheratorStatus": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"totalBytesAreSince": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"totalBytesDropped": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"totalBytesImported": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"totalBytesReturnedToChildren": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"totalBytesReturnedToClients": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"totalBytesReturnedToPeers": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"totalBytesStoredFromOrigin": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"totalBytesStoredFromParents": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"totalBytesStoredFromPeers": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"groupMemberships": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"groupId": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"groupName": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"smartGroup": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// dataSourceJamfProComputerInventoryRead fetches the details of a specific macOS Configuration Profile
// from Jamf Pro using either its unique Name or its ID. The function prioritizes the 'name' attribute over the 'id'
// attribute for fetching details. If neither 'name' nor 'id' is provided, it returns an error.
// Once the details are fetched, they are set in the data source's state.
func dataSourceJamfProComputerInventoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var profile *jamfpro.ResponseComputerInventory
	var err error

	// Fetch profile by 'name' or 'id'
	if v, ok := d.GetOk("name"); ok {
		profileName, ok := v.(string)
		if !ok {
			return diag.Errorf("error asserting 'name' as string")
		}
		profile, err = conn.GetComputerInventoryByName(profileName)
	} else if v, ok := d.GetOk("id"); ok {
		profileID, ok := v.(string)
		if !ok {
			return diag.Errorf("error asserting 'id' as string")
		}
		profile, err = conn.GetComputerInventoryByID(profileID)
	} else {
		return diag.Errorf("Either 'name' or 'id' must be provided")
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to fetch macOS Configuration Profile: %v", err))
	}

	// Set top-level attributes
	d.SetId(profile.ID)
	d.Set("id", profile.ID)
	d.Set("udid", profile.UDID)

	// Set 'general' section
	if err := setGeneralSection(d, profile.General); err != nil {
		return diag.FromErr(err)
	}

	// Set 'diskEncryption' section
	if err := setDiskEncryptionSection(d, profile.DiskEncryption); err != nil {
		return diag.FromErr(err)
	}

	// Set 'purchasing' section
	if err := setPurchasingSection(d, profile.Purchasing); err != nil {
		return diag.FromErr(err)
	}

	// Set 'applications' section
	if err := setApplicationsSection(d, profile.Applications); err != nil {
		return diag.FromErr(err)
	}

	// Set 'storage' section
	if err := setStorageSection(d, profile.Storage); err != nil {
		return diag.FromErr(err)
	}

	// Set 'userAndLocation' section
	if err := setUserAndLocationSection(d, profile.UserAndLocation); err != nil {
		return diag.FromErr(err)
	}

	// Set 'hardware' section
	if err := setHardwareSection(d, profile.Hardware); err != nil {
		return diag.FromErr(err)
	}

	// ... repeat for other sections ...

	return nil
}

// setGeneralSection maps the 'general' section of the computer inventory response to the Terraform schema.
func setGeneralSection(d *schema.ResourceData, general jamfpro.ComputerInventoryDataSubsetGeneral) error {
	// Initialize a map to hold the 'general' section attributes.
	gen := make(map[string]interface{})

	// Map each attribute of the 'general' section from the API response to the corresponding Terraform schema attribute.
	gen["name"] = general.Name
	gen["lastIpAddress"] = general.LastIpAddress
	gen["lastReportedIp"] = general.LastReportedIp
	gen["jamfBinaryVersion"] = general.JamfBinaryVersion
	gen["platform"] = general.Platform
	gen["barcode1"] = general.Barcode1
	gen["barcode2"] = general.Barcode2
	gen["assetTag"] = general.AssetTag
	gen["supervised"] = general.Supervised
	gen["mdmCapable"] = general.MdmCapable.Capable
	gen["reportDate"] = general.ReportDate
	gen["lastContactTime"] = general.LastContactTime
	gen["lastCloudBackupDate"] = general.LastCloudBackupDate
	gen["lastEnrolledDate"] = general.LastEnrolledDate
	gen["mdmProfileExpiration"] = general.MdmProfileExpiration
	gen["initialEntryDate"] = general.InitialEntryDate
	gen["distributionPoint"] = general.DistributionPoint
	gen["itunesStoreAccountActive"] = general.ItunesStoreAccountActive
	gen["enrolledViaAutomatedDeviceEnrollment"] = general.EnrolledViaAutomatedDeviceEnrollment
	gen["userApprovedMdm"] = general.UserApprovedMdm
	gen["declarativeDeviceManagementEnabled"] = general.DeclarativeDeviceManagementEnabled
	gen["managementId"] = general.ManagementId

	// Handle nested object 'remoteManagement'.
	remoteManagement := make(map[string]interface{})
	remoteManagement["managed"] = general.RemoteManagement.Managed
	remoteManagement["managementUsername"] = general.RemoteManagement.ManagementUsername
	gen["remoteManagement"] = []interface{}{remoteManagement}

	// Handle nested object 'site'.
	if general.Site.ID != "" || general.Site.Name != "" {
		site := make(map[string]interface{})
		site["id"] = general.Site.ID
		site["name"] = general.Site.Name
		gen["site"] = []interface{}{site}
	}

	// Handle nested object 'enrollmentMethod'.
	if general.EnrollmentMethod.ID != "" || general.EnrollmentMethod.ObjectName != "" || general.EnrollmentMethod.ObjectType != "" {
		enrollmentMethod := make(map[string]interface{})
		enrollmentMethod["id"] = general.EnrollmentMethod.ID
		enrollmentMethod["objectName"] = general.EnrollmentMethod.ObjectName
		enrollmentMethod["objectType"] = general.EnrollmentMethod.ObjectType
		gen["enrollmentMethod"] = []interface{}{enrollmentMethod}
	}

	// Set the 'general' section in the Terraform resource data.
	return d.Set("general", []interface{}{gen})
}

// setDiskEncryptionSection maps the 'diskEncryption' section of the computer inventory response to the Terraform schema.
func setDiskEncryptionSection(d *schema.ResourceData, diskEncryption jamfpro.ComputerInventoryDataSubsetDiskEncryption) error {
	// Initialize a map to hold the 'diskEncryption' section attributes.
	diskEnc := make(map[string]interface{})

	// Map each attribute of the 'diskEncryption' section from the API response to the corresponding Terraform schema attribute.
	diskEnc["individualRecoveryKeyValidityStatus"] = diskEncryption.IndividualRecoveryKeyValidityStatus
	diskEnc["institutionalRecoveryKeyPresent"] = diskEncryption.InstitutionalRecoveryKeyPresent
	diskEnc["diskEncryptionConfigurationName"] = diskEncryption.DiskEncryptionConfigurationName
	diskEnc["fileVault2EligibilityMessage"] = diskEncryption.FileVault2EligibilityMessage

	// Handle nested object 'bootPartitionEncryptionDetails'.
	bootPartitionDetails := make(map[string]interface{})
	bootPartitionDetails["partitionName"] = diskEncryption.BootPartitionEncryptionDetails.PartitionName
	bootPartitionDetails["partitionFileVault2State"] = diskEncryption.BootPartitionEncryptionDetails.PartitionFileVault2State
	bootPartitionDetails["partitionFileVault2Percent"] = diskEncryption.BootPartitionEncryptionDetails.PartitionFileVault2Percent
	diskEnc["bootPartitionEncryptionDetails"] = []interface{}{bootPartitionDetails}

	// Map 'fileVault2EnabledUserNames' as a list of strings.
	fileVaultUserNames := make([]string, len(diskEncryption.FileVault2EnabledUserNames))
	copy(fileVaultUserNames, diskEncryption.FileVault2EnabledUserNames)

	// Set 'fileVault2EnabledUserNames' in the 'diskEnc' map.
	diskEnc["fileVault2EnabledUserNames"] = fileVaultUserNames

	// Set the 'diskEncryption' section in the Terraform resource data.
	return d.Set("diskEncryption", []interface{}{diskEnc})
}

// setPurchasingSection maps the 'purchasing' section of the computer inventory response to the Terraform schema.
func setPurchasingSection(d *schema.ResourceData, purchasing jamfpro.ComputerInventoryDataSubsetPurchasing) error {
	// Initialize a map to hold the 'purchasing' section attributes.
	purch := make(map[string]interface{})

	// Map each attribute of the 'purchasing' section from the API response to the corresponding Terraform schema attribute.
	purch["leased"] = purchasing.Leased
	purch["purchased"] = purchasing.Purchased
	purch["poNumber"] = purchasing.PoNumber
	purch["poDate"] = purchasing.PoDate
	purch["vendor"] = purchasing.Vendor
	purch["warrantyDate"] = purchasing.WarrantyDate
	purch["appleCareId"] = purchasing.AppleCareId
	purch["leaseDate"] = purchasing.LeaseDate
	purch["purchasePrice"] = purchasing.PurchasePrice
	purch["lifeExpectancy"] = purchasing.LifeExpectancy
	purch["purchasingAccount"] = purchasing.PurchasingAccount
	purch["purchasingContact"] = purchasing.PurchasingContact

	// Map 'extensionAttributes' as a list of maps.
	extAttrs := make([]map[string]interface{}, len(purchasing.ExtensionAttributes))
	for i, attr := range purchasing.ExtensionAttributes {
		attrMap := make(map[string]interface{})
		attrMap["definitionId"] = attr.DefinitionId
		attrMap["name"] = attr.Name
		attrMap["description"] = attr.Description
		attrMap["enabled"] = attr.Enabled
		attrMap["multiValue"] = attr.MultiValue
		attrMap["values"] = attr.Values
		attrMap["dataType"] = attr.DataType
		attrMap["options"] = attr.Options
		attrMap["inputType"] = attr.InputType

		extAttrs[i] = attrMap
	}
	purch["extensionAttributes"] = extAttrs

	// Set the 'purchasing' section in the Terraform resource data.
	return d.Set("purchasing", []interface{}{purch})
}

// setApplicationsSection maps the 'applications' section of the computer inventory response to the Terraform schema.
func setApplicationsSection(d *schema.ResourceData, applications []jamfpro.ComputerInventoryDataSubsetApplication) error {
	// Create a slice to hold the application maps.
	apps := make([]interface{}, len(applications))

	for i, app := range applications {
		// Initialize a map for each application.
		appMap := make(map[string]interface{})

		// Map each attribute of the application from the API response to the corresponding Terraform schema attribute.
		appMap["name"] = app.Name
		appMap["path"] = app.Path
		appMap["version"] = app.Version
		appMap["macAppStore"] = app.MacAppStore
		appMap["sizeMegabytes"] = app.SizeMegabytes
		appMap["bundleId"] = app.BundleId
		appMap["updateAvailable"] = app.UpdateAvailable
		appMap["externalVersionId"] = app.ExternalVersionId

		// Add the application map to the slice.
		apps[i] = appMap
	}

	// Set the 'applications' section in the Terraform resource data.
	return d.Set("applications", apps)
}

// setStorageSection maps the 'storage' section of the computer inventory response to the Terraform schema.
func setStorageSection(d *schema.ResourceData, storage jamfpro.ComputerInventoryDataSubsetStorage) error {
	storageMap := make(map[string]interface{})

	storageMap["bootDriveAvailableSpaceMegabytes"] = storage.BootDriveAvailableSpaceMegabytes

	// Mapping 'disks' array
	disks := make([]interface{}, len(storage.Disks))
	for i, disk := range storage.Disks {
		diskMap := make(map[string]interface{})
		diskMap["id"] = disk.ID
		diskMap["device"] = disk.Device
		diskMap["model"] = disk.Model
		diskMap["revision"] = disk.Revision
		diskMap["serialNumber"] = disk.SerialNumber
		diskMap["sizeMegabytes"] = disk.SizeMegabytes
		diskMap["smartStatus"] = disk.SmartStatus
		diskMap["type"] = disk.Type

		// Map 'partitions' if present
		partitions := make([]interface{}, len(disk.Partitions))
		for j, partition := range disk.Partitions {
			partitionMap := make(map[string]interface{})
			partitionMap["name"] = partition.Name
			partitionMap["sizeMegabytes"] = partition.SizeMegabytes
			partitionMap["availableMegabytes"] = partition.AvailableMegabytes
			partitionMap["partitionType"] = partition.PartitionType
			partitionMap["percentUsed"] = partition.PercentUsed
			partitionMap["fileVault2State"] = partition.FileVault2State
			partitionMap["fileVault2ProgressPercent"] = partition.FileVault2ProgressPercent
			partitionMap["lvmManaged"] = partition.LvmManaged
			partitions[j] = partitionMap
		}
		diskMap["partitions"] = partitions

		disks[i] = diskMap
	}
	storageMap["disks"] = disks

	// Set the 'storage' section in the Terraform resource data.
	return d.Set("storage", []interface{}{storageMap})
}

// setUserAndLocationSection maps the 'userAndLocation' section of the computer inventory response to the Terraform schema.
func setUserAndLocationSection(d *schema.ResourceData, userAndLocation jamfpro.ComputerInventoryDataSubsetUserAndLocation) error {
	userLocationMap := make(map[string]interface{})

	// Map each attribute from the 'userAndLocation' object to the corresponding schema attribute
	userLocationMap["username"] = userAndLocation.Username
	userLocationMap["realname"] = userAndLocation.Realname
	userLocationMap["email"] = userAndLocation.Email
	userLocationMap["position"] = userAndLocation.Position
	userLocationMap["phone"] = userAndLocation.Phone
	userLocationMap["departmentId"] = userAndLocation.DepartmentId
	userLocationMap["buildingId"] = userAndLocation.BuildingId
	userLocationMap["room"] = userAndLocation.Room

	// Map extension attributes if present
	if len(userAndLocation.ExtensionAttributes) > 0 {
		extAttrs := make([]map[string]interface{}, len(userAndLocation.ExtensionAttributes))
		for i, attr := range userAndLocation.ExtensionAttributes {
			attrMap := make(map[string]interface{})
			attrMap["definitionId"] = attr.DefinitionId
			attrMap["name"] = attr.Name
			attrMap["description"] = attr.Description
			attrMap["enabled"] = attr.Enabled
			attrMap["multiValue"] = attr.MultiValue
			attrMap["values"] = attr.Values
			attrMap["dataType"] = attr.DataType
			attrMap["options"] = attr.Options
			attrMap["inputType"] = attr.InputType

			extAttrs[i] = attrMap
		}
		userLocationMap["extensionAttributes"] = extAttrs
	}

	// Set the 'userAndLocation' section in the Terraform resource data
	return d.Set("userAndLocation", []interface{}{userLocationMap})
}

// setHardwareSection maps the 'hardware' section of the computer inventory response to the Terraform schema.
func setHardwareSection(d *schema.ResourceData, hardware jamfpro.ComputerInventoryDataSubsetHardware) error {
	hardwareMap := make(map[string]interface{})

	// Map each attribute from the 'hardware' object to the corresponding schema attribute
	hardwareMap["make"] = hardware.Make
	hardwareMap["model"] = hardware.Model
	hardwareMap["modelIdentifier"] = hardware.ModelIdentifier
	hardwareMap["serialNumber"] = hardware.SerialNumber
	hardwareMap["processorSpeedMhz"] = hardware.ProcessorSpeedMhz
	hardwareMap["processorCount"] = hardware.ProcessorCount
	hardwareMap["coreCount"] = hardware.CoreCount
	hardwareMap["processorType"] = hardware.ProcessorType
	hardwareMap["processorArchitecture"] = hardware.ProcessorArchitecture
	hardwareMap["busSpeedMhz"] = hardware.BusSpeedMhz
	hardwareMap["cacheSizeKilobytes"] = hardware.CacheSizeKilobytes
	hardwareMap["networkAdapterType"] = hardware.NetworkAdapterType
	hardwareMap["macAddress"] = hardware.MacAddress
	hardwareMap["altNetworkAdapterType"] = hardware.AltNetworkAdapterType
	hardwareMap["altMacAddress"] = hardware.AltMacAddress
	hardwareMap["totalRamMegabytes"] = hardware.TotalRamMegabytes
	hardwareMap["openRamSlots"] = hardware.OpenRamSlots
	hardwareMap["batteryCapacityPercent"] = hardware.BatteryCapacityPercent
	hardwareMap["smcVersion"] = hardware.SmcVersion
	hardwareMap["nicSpeed"] = hardware.NicSpeed
	hardwareMap["opticalDrive"] = hardware.OpticalDrive
	hardwareMap["bootRom"] = hardware.BootRom
	hardwareMap["bleCapable"] = hardware.BleCapable
	hardwareMap["supportsIosAppInstalls"] = hardware.SupportsIosAppInstalls
	hardwareMap["appleSilicon"] = hardware.AppleSilicon

	// Map extension attributes if present
	if len(hardware.ExtensionAttributes) > 0 {
		extAttrs := make([]map[string]interface{}, len(hardware.ExtensionAttributes))
		for i, attr := range hardware.ExtensionAttributes {
			attrMap := make(map[string]interface{})
			attrMap["definitionId"] = attr.DefinitionId
			attrMap["name"] = attr.Name
			attrMap["description"] = attr.Description
			attrMap["enabled"] = attr.Enabled
			attrMap["multiValue"] = attr.MultiValue
			attrMap["values"] = attr.Values
			attrMap["dataType"] = attr.DataType
			attrMap["options"] = attr.Options
			attrMap["inputType"] = attr.InputType

			extAttrs[i] = attrMap
		}
		hardwareMap["extensionAttributes"] = extAttrs
	}

	// Set the 'hardware' section in the Terraform resource data
	return d.Set("hardware", []interface{}{hardwareMap})
}

// setHardwareSection maps the 'hardware' section of the computer inventory response to the Terraform schema.
func setLocalUserAccountsSection(d *schema.ResourceData, localUserAccounts []jamfpro.ComputerInventoryDataSubsetLocalUserAccount) error {
	accounts := make([]interface{}, len(localUserAccounts))
	for i, account := range localUserAccounts {
		acc := make(map[string]interface{})
		acc["uid"] = account.UID
		acc["userGuid"] = account.UserGuid
		acc["username"] = account.Username
		acc["fullName"] = account.FullName
		acc["admin"] = account.Admin
		acc["homeDirectory"] = account.HomeDirectory
		acc["homeDirectorySizeMb"] = account.HomeDirectorySizeMb
		acc["fileVault2Enabled"] = account.FileVault2Enabled
		acc["userAccountType"] = account.UserAccountType
		acc["passwordMinLength"] = account.PasswordMinLength
		acc["passwordMaxAge"] = account.PasswordMaxAge
		acc["passwordMinComplexCharacters"] = account.PasswordMinComplexCharacters
		acc["passwordHistoryDepth"] = account.PasswordHistoryDepth
		acc["passwordRequireAlphanumeric"] = account.PasswordRequireAlphanumeric
		acc["computerAzureActiveDirectoryId"] = account.ComputerAzureActiveDirectoryId
		acc["userAzureActiveDirectoryId"] = account.UserAzureActiveDirectoryId
		acc["azureActiveDirectoryId"] = account.AzureActiveDirectoryId
		accounts[i] = acc
	}
	return d.Set("localUserAccounts", accounts)
}

// setCertificatesSection maps the 'certificate' section of the computer inventory response to the Terraform schema.
func setCertificatesSection(d *schema.ResourceData, certificates []jamfpro.ComputerInventoryDataSubsetCertificate) error {
	certs := make([]interface{}, len(certificates))
	for i, cert := range certificates {
		certMap := make(map[string]interface{})
		certMap["commonName"] = cert.CommonName
		certMap["identity"] = cert.Identity
		certMap["expirationDate"] = cert.ExpirationDate
		certMap["username"] = cert.Username
		certMap["lifecycleStatus"] = cert.LifecycleStatus
		certMap["certificateStatus"] = cert.CertificateStatus
		certMap["subjectName"] = cert.SubjectName
		certMap["serialNumber"] = cert.SerialNumber
		certMap["sha1Fingerprint"] = cert.Sha1Fingerprint
		certMap["issuedDate"] = cert.IssuedDate
		certs[i] = certMap
	}
	return d.Set("certificates", certs)
}

// setAttachmentsSection maps the 'attachments' section of the computer inventory response to the Terraform schema.
func setAttachmentsSection(d *schema.ResourceData, attachments []jamfpro.ComputerInventoryDataSubsetAttachment) error {
	atts := make([]interface{}, len(attachments))
	for i, att := range attachments {
		attMap := make(map[string]interface{})
		attMap["id"] = att.ID
		attMap["name"] = att.Name
		attMap["fileType"] = att.FileType
		attMap["sizeBytes"] = att.SizeBytes
		atts[i] = attMap
	}
	return d.Set("attachments", atts)
}

// setPluginsSection maps the 'plugins' section of the computer inventory response to the Terraform schema.
func setPluginsSection(d *schema.ResourceData, plugins []jamfpro.ComputerInventoryDataSubsetPlugin) error {
	pluginList := make([]interface{}, len(plugins))
	for i, plugin := range plugins {
		pluginMap := make(map[string]interface{})
		pluginMap["name"] = plugin.Name
		pluginMap["version"] = plugin.Version
		pluginMap["path"] = plugin.Path
		pluginList[i] = pluginMap
	}
	return d.Set("plugins", pluginList)
}

// setPackageReceiptsSection maps the 'package receipts' section of the computer inventory response to the Terraform schema.
func setPackageReceiptsSection(d *schema.ResourceData, packageReceipts jamfpro.ComputerInventoryDataSubsetPackageReceipts) error {
	prMap := make(map[string]interface{})
	prMap["installedByJamfPro"] = packageReceipts.InstalledByJamfPro
	prMap["installedByInstallerSwu"] = packageReceipts.InstalledByInstallerSwu
	prMap["cached"] = packageReceipts.Cached
	return d.Set("packageReceipts", []interface{}{prMap})
}

func setFontsSection(d *schema.ResourceData, fonts []jamfpro.ComputerInventoryDataSubsetFont) error {
	fontsList := make([]interface{}, len(fonts))
	for i, font := range fonts {
		fontMap := make(map[string]interface{})
		fontMap["name"] = font.Name
		fontMap["version"] = font.Version
		fontMap["path"] = font.Path
		fontsList[i] = fontMap
	}
	return d.Set("fonts", fontsList)
}

func setSecuritySection(d *schema.ResourceData, security jamfpro.ComputerInventoryDataSubsetSecurity) error {
	securityMap := make(map[string]interface{})
	securityMap["sipStatus"] = security.SipStatus
	securityMap["gatekeeperStatus"] = security.GatekeeperStatus
	securityMap["xprotectVersion"] = security.XprotectVersion
	securityMap["autoLoginDisabled"] = security.AutoLoginDisabled
	securityMap["remoteDesktopEnabled"] = security.RemoteDesktopEnabled
	securityMap["activationLockEnabled"] = security.ActivationLockEnabled
	securityMap["recoveryLockEnabled"] = security.RecoveryLockEnabled
	securityMap["firewallEnabled"] = security.FirewallEnabled
	securityMap["secureBootLevel"] = security.SecureBootLevel
	securityMap["externalBootLevel"] = security.ExternalBootLevel
	securityMap["bootstrapTokenAllowed"] = security.BootstrapTokenAllowed
	return d.Set("security", []interface{}{securityMap})
}

func setOperatingSystemSection(d *schema.ResourceData, operatingSystem jamfpro.ComputerInventoryDataSubsetOperatingSystem) error {
	osMap := make(map[string]interface{})
	osMap["name"] = operatingSystem.Name
	osMap["version"] = operatingSystem.Version
	osMap["build"] = operatingSystem.Build
	osMap["supplementalBuildVersion"] = operatingSystem.SupplementalBuildVersion
	osMap["rapidSecurityResponse"] = operatingSystem.RapidSecurityResponse
	osMap["activeDirectoryStatus"] = operatingSystem.ActiveDirectoryStatus
	osMap["fileVault2Status"] = operatingSystem.FileVault2Status
	osMap["softwareUpdateDeviceId"] = operatingSystem.SoftwareUpdateDeviceId
	// Map extension attributes if present
	extAttrs := make([]map[string]interface{}, len(operatingSystem.ExtensionAttributes))
	for i, attr := range operatingSystem.ExtensionAttributes {
		attrMap := make(map[string]interface{})
		attrMap["definitionId"] = attr.DefinitionId
		attrMap["name"] = attr.Name
		attrMap["description"] = attr.Description
		attrMap["enabled"] = attr.Enabled
		attrMap["multiValue"] = attr.MultiValue
		attrMap["values"] = attr.Values
		attrMap["dataType"] = attr.DataType
		attrMap["options"] = attr.Options
		attrMap["inputType"] = attr.InputType

		extAttrs[i] = attrMap
	}
	osMap["extensionAttributes"] = extAttrs
	return d.Set("operatingSystem", []interface{}{osMap})
}

func setLicensedSoftwareSection(d *schema.ResourceData, licensedSoftware []jamfpro.ComputerInventoryDataSubsetLicensedSoftware) error {
	softwareList := make([]interface{}, len(licensedSoftware))
	for i, software := range licensedSoftware {
		softwareMap := make(map[string]interface{})
		softwareMap["id"] = software.ID
		softwareMap["name"] = software.Name
		softwareList[i] = softwareMap
	}
	return d.Set("licensedSoftware", softwareList)
}

func setIBeaconsSection(d *schema.ResourceData, ibeacons []jamfpro.ComputerInventoryDataSubsetIbeacon) error {
	ibeaconList := make([]interface{}, len(ibeacons))
	for i, ibeacon := range ibeacons {
		ibeaconMap := make(map[string]interface{})
		ibeaconMap["name"] = ibeacon.Name
		ibeaconList[i] = ibeaconMap
	}
	return d.Set("ibeacons", ibeaconList)
}

func setSoftwareUpdatesSection(d *schema.ResourceData, softwareUpdates []jamfpro.ComputerInventoryDataSubsetSoftwareUpdate) error {
	updateList := make([]interface{}, len(softwareUpdates))
	for i, update := range softwareUpdates {
		updateMap := make(map[string]interface{})
		updateMap["name"] = update.Name
		updateMap["version"] = update.Version
		updateMap["packageName"] = update.PackageName
		updateList[i] = updateMap
	}
	return d.Set("softwareUpdates", updateList)
}

func setExtensionAttributesSection(d *schema.ResourceData, extensionAttributes []jamfpro.ComputerInventoryDataSubsetExtensionAttribute) error {
	attrList := make([]interface{}, len(extensionAttributes))
	for i, attr := range extensionAttributes {
		attrMap := make(map[string]interface{})
		attrMap["definitionId"] = attr.DefinitionId
		attrMap["name"] = attr.Name
		attrMap["description"] = attr.Description
		attrMap["enabled"] = attr.Enabled
		attrMap["multiValue"] = attr.MultiValue
		attrMap["values"] = attr.Values
		attrMap["dataType"] = attr.DataType
		attrMap["options"] = attr.Options
		attrMap["inputType"] = attr.InputType
		attrList[i] = attrMap
	}
	return d.Set("extensionAttributes", attrList)
}

func setContentCachingSection(d *schema.ResourceData, contentCaching jamfpro.ComputerInventoryDataSubsetContentCaching) error {
	cachingMap := make(map[string]interface{})
	cachingMap["computerContentCachingInformationId"] = contentCaching.ComputerContentCachingInformationId
	// ... map other content caching attributes ...

	// Handle nested objects like 'parents', 'alerts', etc. if needed
	// ...

	return d.Set("contentCaching", []interface{}{cachingMap})
}

func setAlertsSection(d *schema.ResourceData, alerts []jamfpro.ComputerInventoryDataSubsetContentCachingAlert) error {
	alertsList := make([]interface{}, len(alerts))
	for i, alert := range alerts {
		alertMap := make(map[string]interface{})
		alertMap["contentCachingParentAlertId"] = alert.ContentCachingParentAlertId
		alertMap["addresses"] = alert.Addresses
		alertMap["className"] = alert.ClassName
		alertMap["postDate"] = alert.PostDate
		// ... map other alert attributes ...
		alertsList[i] = alertMap
	}
	return d.Set("alerts", alertsList)
}
