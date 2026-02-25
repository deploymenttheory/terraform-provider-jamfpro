package computer_inventory

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProMacOSComputerInventory provides information about a specific computer's inventory by its ID or Name.
func DataSourceJamfProComputerInventory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"warn_if_not_found": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enabling this setting will cause the provider to only WARN if a computer is not found. By default the provider will ERROR.",
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"serial_number": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"udid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"general": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_reported_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"jamf_binary_version": {
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
						"asset_tag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_management": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"managed": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"management_username": {
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
						"mdm_capable": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"capable": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"capable_users": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"report_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_contact_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_cloud_backup_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_enrolled_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mdm_profile_expiration": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"initial_entry_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"distribution_point": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enrollment_method": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"object_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"object_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"site_id": {
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
						"itunes_store_account_active": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enrolled_via_automated_device_enrollment": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"user_approved_mdm": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"declarative_device_management_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"extension_attributes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition_id": {
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
									"multi_value": {
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
									"data_type": {
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
									"input_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"management_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"disk_encryption": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"boot_partition_encryption_details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"partition_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"partition_file_vault2_state": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"partition_file_vault2_percent": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"individual_recovery_key_validity_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"institutional_recovery_key_present": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"disk_encryption_configuration_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_vault2_enabled_user_names": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"file_vault2_eligibility_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"purchasing": {
				Type:     schema.TypeList,
				Computed: true,
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
						"po_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"po_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vendor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"warranty_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"apple_care_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lease_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"purchase_price": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"life_expectancy": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"purchasing_account": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"purchasing_contact": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extension_attributes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition_id": {
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
									"multi_value": {
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
									"data_type": {
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
									"input_type": {
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
						"mac_app_store": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"size_megabytes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bundle_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_available": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"external_version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"storage": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"boot_drive_available_space_megabytes": {
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
									"serial_number": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size_megabytes": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"smart_status": {
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
												"size_megabytes": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"available_megabytes": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"partition_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"percent_used": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"file_vault2_state": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"file_vault2_progress_percent": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"lvm_managed": {
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
			"user_and_location": {
				Type:     schema.TypeList,
				Computed: true,
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
						"department_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"building_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"room": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extension_attributes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition_id": {
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
									"multi_value": {
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
									"data_type": {
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
									"input_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"configuration_profiles": {
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
						"last_installed": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"removable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"profile_identifier": {
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
						"model_identifier": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"serial_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"processor_speed_mhz": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"processor_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"core_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"processor_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"processor_architecture": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bus_speed_mhz": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cache_size_kilobytes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"network_adapter_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alt_network_adapter_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alt_mac_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"total_ram_megabytes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"open_ram_slots": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"battery_capacity_percent": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"smc_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nic_speed": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"optical_drive": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"boot_rom": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ble_capable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"supports_ios_app_installs": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"apple_silicon": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"extension_attributes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition_id": {
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
									"multi_value": {
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
									"data_type": {
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
									"input_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"local_user_accounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_guid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"full_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"admin": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"home_directory": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"home_directory_size_mb": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"file_vault2_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"user_account_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password_min_length": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"password_max_age": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"password_min_complex_characters": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"password_history_depth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"password_require_alphanumeric": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"computer_azure_active_directory_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_azure_active_directory_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"azure_active_directory_id": {
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
						"common_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"identity": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"expiration_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lifecycle_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subject_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"serial_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sha1_fingerprint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"issued_date": {
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
						"file_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size_bytes": {
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
			"package_receipts": {
				Type:     schema.TypeList,
				Computed: true,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"installed_by_jamf_pro": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"installed_by_installer_swu": {
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

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sip_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gatekeeper_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"xprotect_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_login_disabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"remote_desktop_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"activation_lock_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"recovery_lock_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"firewall_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"secure_boot_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_boot_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bootstrap_token_allowed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"operating_system": {
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
						"build": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"supplemental_build_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rapid_security_response": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"active_directory_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filevault2_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"software_update_device_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extension_attributes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition_id": {
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
									"multi_value": {
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
									"data_type": {
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
									"input_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"licensed_software": {
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
			"software_updates": {
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
						"package_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"extension_attributes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition_id": {
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
						"multi_value": {
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
						"data_type": {
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
						"input_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			// ignoring contentCashing section intentionally as it doesnt serve any use in terraform scenarios
			"group_memberships": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"smart_group": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
