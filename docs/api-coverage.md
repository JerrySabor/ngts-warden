# NGTS API coverage

Generated from `api/openapi.json`. This snapshot contains **147 operations** across **24 resource groups**.

| CLI command | Method | Path | Operation ID | Required inputs | Body | Response content types | Summary |
|---|---|---|---|---|---|---|---|
| `built-in-accounts create-v1-serviceaccounts` | `POST` | `/v1/serviceaccounts` | `create-v1-serviceaccounts` | `—` | `required` | `application/json` | Creates a Service Account |
| `built-in-accounts delete-v1-serviceaccounts-by-id` | `DELETE` | `/v1/serviceaccounts/{id}` | `delete-v1-serviceaccounts-byId` | `path:id` | `—` | `—` | Deletes a Service Account |
| `built-in-accounts get-v1-serviceaccounts` | `GET` | `/v1/serviceaccounts` | `get-v1-serviceaccounts` | `—` | `—` | `application/json` | Retrieves all the Service Accounts the |
| `built-in-accounts get-v1-serviceaccounts-by-id` | `GET` | `/v1/serviceaccounts/{id}` | `get-v1-serviceaccounts-byId` | `path:id` | `—` | `application/json` | Gets a Service Account |
| `built-in-accounts get-v1-serviceaccountscopes` | `GET` | `/v1/serviceaccounts/scopes` | `get-v1-serviceaccountscopes` | `—` | `—` | `application/json` | Retrieves all the Service Accounts Scopes |
| `built-in-accounts patch-v1-serviceaccounts-by-id` | `PATCH` | `/v1/serviceaccounts/{id}` | `patch-v1-serviceaccounts-byId` | `path:id` | `required` | `—` | Updates a Service Account |
| `built-in-accounts put-v1-serviceaccounts-by-id-credentials` | `PUT` | `/v1/serviceaccounts/{id}/credentials` | `put-v1-serviceaccounts-byId-credentials` | `path:id` | `required` | `application/json` | Updates a Service Account credentials |
| `built-in-accounts put-v1-serviceaccounts-by-id-ocitoken` | `PUT` | `/v1/serviceaccounts/{id}/ocitoken` | `put-v1-serviceaccounts-byId-ocitoken` | `path:id` | `—` | `application/json` | Regenerate the OCI registry token for |
| `certificate-approvals certificaterequests-approval-rule-create` | `POST` | `/v1/certificaterequests/approvalrules` | `certificaterequests_approval_rule_create` | `—` | `required` | `application/json` | Create an approval rule for certificate |
| `certificate-approvals certificaterequests-approval-rule-delete` | `DELETE` | `/v1/certificaterequests/approvalrules/{id}` | `certificaterequests_approval_rule_delete` | `path:id` | `—` | `application/json` | Delete certificate request workflow approval rule |
| `certificate-approvals certificaterequests-approval-rule-get-by-id` | `GET` | `/v1/certificaterequests/approvalrules/{id}` | `certificaterequests_approval_rule_getById` | `path:id` | `—` | `application/json` | Retrieve approval rule by id |
| `certificate-approvals certificaterequests-approval-rule-update` | `PUT` | `/v1/certificaterequests/approvalrules/{id}` | `certificaterequests_approval_rule_update` | `path:id` | `optional` | `application/json` | Update certificate request workflow approval rule |
| `certificate-approvals certificaterequests-approval-rules-get-all` | `GET` | `/v1/certificaterequests/approvalrules` | `certificaterequests_approval_rules_getAll` | `—` | `—` | `application/json` | Get all approval rules |
| `certificate-approvals certificaterequests-approvalrequest` | `GET` | `/v1/certificaterequests/approvalrequests/{entityId}` | `certificaterequests_approvalrequest` | `path:entityId` | `—` | `application/json` | Retrieve approval request for specific certificate |
| `certificate-approvals certificaterequests-approve` | `POST` | `/v1/certificaterequests/{id}/approval/{decision}` | `certificaterequests_approve` | `path:id, path:decision` | `optional` | `application/json` | Approve or reject pending certificate request |
| `certificate-approvals certificaterequests-bulk-approve` | `POST` | `/v1/certificaterequests/approval/bulk/{decision}` | `certificaterequests_bulk_approve` | `path:decision` | `optional` | `application/json` | Approve or reject multiple pending approval |
| `certificate-discovery integrationsservices-create` | `POST` | `/v1/integrationservices` | `integrationsservices_create` | `—` | `required` | `application/json` | Add a service |
| `certificate-discovery integrationsservices-delete` | `DELETE` | `/v1/integrationservices/{id}` | `integrationsservices_delete` | `path:id` | `—` | `application/json` | Remove a service |
| `certificate-discovery integrationsservices-get-all` | `GET` | `/v1/integrationservices` | `integrationsservices_getAll` | `—` | `—` | `application/json` | Get a list of services |
| `certificate-discovery integrationsservices-get-by-id` | `GET` | `/v1/integrationservices/{id}` | `integrationsservices_getById` | `path:id` | `—` | `application/json` | Get service details |
| `certificate-discovery integrationsservices-update` | `PATCH` | `/v1/integrationservices/{id}` | `integrationsservices_update` | `path:id` | `required` | `application/json` | Update Service properties |
| `certificate-import certificateimports-create` | `POST` | `/outagedetection/v1/certificates` | `certificateimports_create` | `—` | `optional` | `application/json` | Import a set of raw certificates |
| `certificate-policies certificateissuingtemplate-create` | `POST` | `/v1/certificateissuingtemplates` | `certificateissuingtemplate_create` | `—` | `optional` | `application/json` | Add an issuing template |
| `certificate-policies certificateissuingtemplate-delete` | `DELETE` | `/v1/certificateissuingtemplates/{id}` | `certificateissuingtemplate_delete` | `path:id` | `—` | `application/json` | Remove an issuing template |
| `certificate-policies certificateissuingtemplate-get-all` | `GET` | `/v1/certificateissuingtemplates` | `certificateissuingtemplate_getAll` | `—` | `—` | `application/json` | Get the details of issuing templates |
| `certificate-policies certificateissuingtemplate-get-by-id` | `GET` | `/v1/certificateissuingtemplates/{id}` | `certificateissuingtemplate_getById` | `path:id` | `—` | `application/json` | Get an issuing template details |
| `certificate-policies certificateissuingtemplate-update` | `PUT` | `/v1/certificateissuingtemplates/{id}` | `certificateissuingtemplate_update` | `path:id` | `optional` | `application/json` | Overwrite an issuing template details |
| `certificate-policies domainssynchronization` | `POST` | `/v1/certificateissuingtemplates/domainssynchronization` | `domainssynchronization` | `—` | `optional` | `application/json` | Synchronize issuing templates domains with CA |
| `certificate-requests certificaterequests-create` | `POST` | `/outagedetection/v1/certificaterequests` | `certificaterequests_create` | `—` | `optional` | `application/json` | Create a certificate request |
| `certificate-requests certificaterequests-get-all` | `GET` | `/outagedetection/v1/certificaterequests` | `certificaterequests_getAll` | `—` | `—` | `application/json` | Get the details of all certificate |
| `certificate-requests certificaterequests-get-by-id` | `GET` | `/outagedetection/v1/certificaterequests/{id}` | `certificaterequests_getById` | `path:id` | `—` | `application/json` | Get a certificate request details |
| `certificate-requests certificaterequests-resubmit-by-id` | `POST` | `/outagedetection/v1/certificaterequests/{id}/resubmission` | `certificaterequests_resubmitById` | `path:id` | `optional` | `application/json` | Resubmit a certificate request |
| `certificate-requests certificaterequests-validation` | `POST` | `/outagedetection/v1/certificaterequests/validation` | `certificaterequests_validation` | `—` | `optional` | `application/json` | Validate a certificate request |
| `certificate-requests get-certificate-requests-by-expression` | `POST` | `/outagedetection/v1/certificaterequestssearch` | `getCertificateRequestsByExpression` | `—` | `optional` | `application/json` | Get the details of certificate requests |
| `certificate-tags tag-values-create` | `POST` | `/v1/tags/{name}/values` | `tag_values_create` | `path:name` | `required` | `application/json` | Create tag values |
| `certificate-tags tags-assign-to-entities` | `PATCH` | `/v1/tagsassignment` | `tags_assignToEntities` | `—` | `required` | `application/json` | Replace Add Or Delete Tags |
| `certificate-tags tags-assignment-aggregates` | `POST` | `/v1/tagsassignment/aggregates` | `tags_assignmentAggregates` | `—` | `required` | `application/json` | Bulk operation to retrieve number of |
| `certificate-tags tags-bulk-create` | `POST` | `/v1/tags/creation` | `tags_bulk_create` | `—` | `required` | `application/json` | Create tags in bulk |
| `certificate-tags tags-bulk-delete` | `POST` | `/v1/tags/deletion` | `tags_bulk_delete` | `—` | `required` | `application/json` | Delete tags in bulk |
| `certificate-tags tags-create` | `POST` | `/v1/tags` | `tags_create` | `—` | `required` | `application/json` | Create a tag |
| `certificate-tags tags-delete-by-name` | `DELETE` | `/v1/tags/{name}` | `tags_deleteByName` | `path:name` | `—` | `application/json` | Delete tag by name |
| `certificate-tags tags-delete-value-by-name` | `DELETE` | `/v1/tags/{name}/values/{value}` | `tags_deleteValueByName` | `path:name, path:value` | `—` | `application/json` | Delete a tag value |
| `certificate-tags tags-get-all` | `GET` | `/v1/tags` | `tags_getAll` | `—` | `—` | `application/json` | Retrieve all tags |
| `certificate-tags tags-get-all-values` | `GET` | `/v1/tags/values` | `tags_getAllValues` | `—` | `—` | `application/json` | Retrieve values for all tags |
| `certificate-tags tags-get-by-name` | `GET` | `/v1/tags/{name}` | `tags_getByName` | `path:name` | `—` | `application/json` | Retrieve tag by name |
| `certificate-tags tags-get-values` | `GET` | `/v1/tags/{name}/values` | `tags_get_values` | `path:name` | `—` | `application/json` | Retrieve values for a tag |
| `certificates certificateretirement-delete-certificates` | `POST` | `/outagedetection/v1/certificates/deletion` | `certificateretirement_deleteCertificates` | `—` | `required` | `application/json` | Delete a set of retired certificates |
| `certificates certificateretirement-recover-certificates` | `POST` | `/outagedetection/v1/certificates/recovery` | `certificateretirement_recoverCertificates` | `—` | `optional` | `application/json` | Recover a set of certificates |
| `certificates certificateretirement-retire-certificates` | `POST` | `/outagedetection/v1/certificates/retirement` | `certificateretirement_retireCertificates` | `—` | `optional` | `application/json` | Retire certificates |
| `certificates get-all` | `GET` | `/outagedetection/v1/certificates` | `certificates_getAll` | `—` | `—` | `application/json, text/csv` | Retrieve all certificate data |
| `certificates get-by-id` | `GET` | `/outagedetection/v1/certificates/{id}` | `certificates_getById` | `path:id` | `—` | `application/json` | Get a certificate details |
| `certificates get-contents-by-id` | `GET` | `/outagedetection/v1/certificates/{id}/contents` | `certificates_getContentsById` | `path:id` | `—` | `application/json, application/octet-stream, text/plain` | Download a certificate |
| `certificates search-get-by-expression` | `POST` | `/outagedetection/v1/certificatesearch` | `certificates_search_getByExpression` | `—` | `optional` | `application/json, text/csv` | Retrieve certificate data matching search criteria |
| `certificates validation` | `POST` | `/outagedetection/v1/certificates/validation` | `certificates_validation` | `—` | `optional` | `application/json` | Request validation for a set of |
| `credentials delete-public-cms-conf-id` | `DELETE` | `/v1/credentialmanagerconfigurations/{id}` | `delete-public-cms-conf-id` | `path:id` | `—` | `application/json` | Delete a Credential Manager Service configuration |
| `credentials delete-public-cms-credential` | `DELETE` | `/v1/credentials` | `delete-public-cms-credential` | `—` | `—` | `application/json` | Delete shared credentials |
| `credentials delete-public-cms-credential-id` | `DELETE` | `/v1/credentials/{id}` | `delete-public-cms-credential-id` | `path:id` | `—` | `application/json` | Delete shared credential by ID |
| `credentials get-public-cms-conf` | `GET` | `/v1/credentialmanagerconfigurations` | `get-public-cms-conf` | `—` | `—` | `application/json` | Retrieves a set of Credential Manager |
| `credentials get-public-cms-conf-id` | `GET` | `/v1/credentialmanagerconfigurations/{id}` | `get-public-cms-conf-id` | `path:id` | `—` | `application/json` | Retrieves a Credential Manager Service configurati |
| `credentials get-public-cms-credential` | `GET` | `/v1/credentials` | `get-public-cms-credential` | `—` | `—` | `application/json` | Retrieves credentials for a company |
| `credentials get-public-cms-credential-id` | `GET` | `/v1/credentials/{id}` | `get-public-cms-credential-id` | `path:id` | `—` | `application/json` | Retrieves shared credential by ID |
| `credentials post-public-cms-conf` | `POST` | `/v1/credentialmanagerconfigurations` | `post-public-cms-conf` | `—` | `optional` | `application/json` | Add a set of Credential Manager |
| `credentials post-public-cms-conf-test` | `POST` | `/v1/credentialmanagerconfigurations/test` | `post-public-cms-conf-test` | `—` | `optional` | `application/json` | Test the connection to a privileged |
| `credentials post-public-cms-conf-test-id` | `POST` | `/v1/credentialmanagerconfigurations/{id}/test` | `post-public-cms-conf-test-id` | `path:id` | `—` | `application/json` | Test the connection to an external |
| `credentials post-public-cms-credential` | `POST` | `/v1/credentials` | `post-public-cms-credential` | `—` | `optional` | `application/json` | Add a set of new shared |
| `credentials post-public-cms-credential-test-id` | `POST` | `/v1/credentials/test` | `post-public-cms-credential-test-id` | `—` | `optional` | `application/json` | Test the access to shared credential |
| `credentials put-public-cms-conf` | `PUT` | `/v1/credentialmanagerconfigurations` | `put-public-cms-conf` | `—` | `optional` | `application/json` | Update a Credential Manager Service configuration |
| `credentials put-public-cms-credential` | `PUT` | `/v1/credentials` | `put-public-cms-credential` | `—` | `optional` | `application/json` | Update a shared credential |
| `event-logs activitylogs-get-all-by-expression` | `POST` | `/v1/activitylogsearch/export` | `activitylogs_getAllByExpression` | `—` | `optional` | `application/json, text/csv` | Export filtered event log data to |
| `event-logs activitylogs-get-by-expression` | `POST` | `/v1/activitylogsearch` | `activitylogs_getByExpression` | `—` | `optional` | `application/json` | Retrieve count and activity log entries |
| `event-logs activitylogtypes-get` | `GET` | `/v1/activitytypes` | `activitylogtypes_get` | `—` | `—` | `application/json` | Retrieve types of activities used for |
| `inventory-monitoring get-v1-tenant-expiration-notification-configuration` | `GET` | `/v1/expirationnotifications/tenantconfiguration` | `get-v1-tenant-expiration-notification-configuration` | `—` | `—` | `application/json` | Retrieve the certificate expiration notification c |
| `inventory-monitoring inventorymonitoringconfiguration-get-by-type` | `GET` | `/outagedetection/v1/inventorymonitoringconfig/{type}` | `inventorymonitoringconfiguration_getByType` | `path:type` | `—` | `application/json` | Get the details of the current |
| `inventory-monitoring inventorymonitoringconfiguration-update` | `PUT` | `/outagedetection/v1/inventorymonitoringconfig/{type}` | `inventorymonitoringconfiguration_update` | `path:type` | `required` | `application/json` | Updates existing inventory monitoring configuratio |
| `inventory-monitoring inventorymonitoringconfigurationscheduler-update` | `PUT` | `/outagedetection/v1/inventorymonitoringconfig/{type}/scheduler` | `inventorymonitoringconfigurationscheduler_update` | `path:type` | `—` | `application/json` | Update inventory monitoring scheduler by type |
| `inventory-monitoring put-v1-tenant-expiration-notification-configuration` | `PUT` | `/v1/expirationnotifications/tenantconfiguration` | `put-v1-tenant-expiration-notification-configuration` | `—` | `required` | `application/json` | Update the certificate expiration notification con |
| `issuer-certificates intermediatecertificates-get-all` | `GET` | `/v1/distributedissuers/intermediatecertificates` | `intermediatecertificates_getAll` | `—` | `—` | `application/json` | Get the details of all Issuer |
| `issuer-configurations configurations-create` | `POST` | `/v1/distributedissuers/configurations` | `configurations_create` | `—` | `optional` | `application/json` | Create a new Issuer configuration |
| `issuer-configurations configurations-delete` | `DELETE` | `/v1/distributedissuers/configurations/{id}` | `configurations_delete` | `path:id` | `—` | `application/json` | Remove an Issuer configuration |
| `issuer-configurations configurations-get-all` | `GET` | `/v1/distributedissuers/configurations` | `configurations_getAll` | `—` | `—` | `application/json` | Get the details of all Issuer |
| `issuer-configurations configurations-get-by-id` | `GET` | `/v1/distributedissuers/configurations/{id}` | `configurations_getById` | `path:id` | `—` | `application/json` | Get configurations details for a specific |
| `issuer-configurations configurations-update` | `PATCH` | `/v1/distributedissuers/configurations/{id}` | `configurations_update` | `path:id` | `required` | `application/json` | Update an Issuer configuration details |
| `machine-installations get-machine-identities-by-expression` | `POST` | `/v1/machineidentitysearch` | `getMachineIdentitiesByExpression` | `—` | `optional` | `application/json` | Get the details of machine identities |
| `machine-installations machineidentities-create` | `POST` | `/v1/machineidentities` | `machineidentities_create` | `—` | `optional` | `application/json` | Add a machine identity to a |
| `machine-installations machineidentities-delete` | `DELETE` | `/v1/machineidentities/{id}` | `machineidentities_delete` | `path:id` | `—` | `application/json` | Remove a machine identity |
| `machine-installations machineidentities-get-all` | `GET` | `/v1/machineidentities` | `machineidentities_getAll` | `—` | `—` | `application/json` | Get the details of all machine |
| `machine-installations machineidentities-get-by-id` | `GET` | `/v1/machineidentities/{id}` | `machineidentities_getById` | `path:id` | `—` | `application/json` | Get a machine identity details |
| `machine-installations machineidentities-initiate-workflow` | `POST` | `/v1/machineidentities/{id}/workflows` | `machineidentities_initiateWorkflow` | `path:id` | `optional` | `application/json` | Initiate a machine workflow |
| `machine-installations machineidentities-update` | `PATCH` | `/v1/machineidentities/{id}` | `machineidentities_update` | `path:id` | `optional` | `application/json` | Update a machine identity details |
| `machine-types get-all` | `GET` | `/v1/machinetypes` | `machineTypes_getAll` | `—` | `—` | `application/json` | List Machine Types |
| `machines abort-v1-batchprovisionings-for-machine-id` | `POST` | `/v1/machines/{id}/batchprovisionings/abort` | `abort-v1-batchprovisionings-forMachineId` | `path:id` | `—` | `application/json` | Abort active batch provisioning for a |
| `machines create` | `POST` | `/v1/machines` | `machines_create` | `—` | `optional` | `application/json` | Add a machine |
| `machines delete` | `DELETE` | `/v1/machines/{id}` | `machines_delete` | `path:id` | `—` | `application/json` | Delete a machine |
| `machines discovery-results-abortdiscovery` | `POST` | `/v1/machines/{id}/discovery/abort` | `machineDiscoveryResults_abortdiscovery` | `path:id` | `—` | `application/json` | Abort machine discovery |
| `machines discovery-results-get-by-machine-id` | `GET` | `/v1/machines/{id}/discovery` | `machineDiscoveryResults_getByMachineId` | `path:id` | `—` | `application/json` | Get the discovery results for a |
| `machines get-all` | `GET` | `/v1/machines` | `machines_getAll` | `—` | `—` | `application/json` | Get the details of all machines |
| `machines get-by-id` | `GET` | `/v1/machines/{id}` | `machines_getById` | `path:id` | `—` | `application/json` | Get a machine details |
| `machines get-machines-by-expression` | `POST` | `/v1/machinesearch` | `getMachinesByExpression` | `—` | `optional` | `application/json` | Get the details of machines matching |
| `machines initiate-workflow` | `POST` | `/v1/machines/{id}/workflows` | `machines_initiateWorkflow` | `path:id` | `optional` | `application/json` | Initiate the workflow |
| `machines update` | `PATCH` | `/v1/machines/{id}` | `machines_update` | `path:id` | `optional` | `application/json` | Update a machine details |
| `plugins delete-v1-plugins-id` | `DELETE` | `/v1/plugins/{id}` | `delete-v1-plugins-id` | `path:id` | `—` | `—` | Delete a local plugin |
| `plugins delete-v1-plugins-id-exclusions` | `DELETE` | `/v1/plugins/{id}/disablements` | `delete-v1-plugins-id-exclusions` | `path:id` | `—` | `—` | Remove plugin disablement |
| `plugins get-v1-plugins` | `GET` | `/v1/plugins` | `get-v1-plugins` | `—` | `—` | `application/json` | Retrieve all plugins |
| `plugins get-v1-plugins-exclusions` | `GET` | `/v1/plugins/disablements` | `get-v1-plugins-exclusions` | `—` | `—` | `application/json` | Retrieve all disabled plugins |
| `plugins get-v1-plugins-id` | `GET` | `/v1/plugins/{id}` | `get-v1-plugins-id` | `path:id` | `—` | `application/json` | Retrieve plugin by ID |
| `plugins patch-v1-plugins-id` | `PATCH` | `/v1/plugins/{id}` | `patch-v1-plugins-id` | `path:id` | `optional` | `—` | Update a local plugin |
| `plugins post-v1-plugins` | `POST` | `/v1/plugins` | `post-v1-plugins` | `—` | `optional` | `application/json` | Create a local plugin |
| `plugins post-v1-plugins-id-exclusions` | `POST` | `/v1/plugins/{id}/disablements` | `post-v1-plugins-id-exclusions` | `path:id` | `—` | `application/json` | Disable a plugin |
| `private-key-import certificates-import` | `POST` | `/v1/certificates/imports` | `certificates_import` | `—` | `optional` | `application/json` | Import a list of certificates and |
| `private-key-import certificates-import-get-by-import-id` | `GET` | `/v1/certificates/imports/{id}` | `certificatesImport_getByImportId` | `path:id` | `—` | `application/json` | Retrieve import details |
| `renewal-monitoring get-v1-status` | `GET` | `/v1/autorenewal/status` | `get-v1-status` | `—` | `—` | `application/json` | Get the current certificate auto-renewal monitorin |
| `renewal-monitoring get-v1-tenant-renewal-configuration` | `GET` | `/v1/autorenewal/tenantconfiguration` | `get-v1-tenant-renewal-configuration` | `—` | `—` | `application/json` | Retrieve the monitoring configuration |
| `renewal-monitoring post-v1-run-autorenewal` | `POST` | `/v1/autorenewal/trigger` | `post-v1-run-autorenewal` | `—` | `—` | `application/json` | Attempt to initiate the certificate renewal |
| `renewal-monitoring put-v1-tenant-renewal-configuration` | `PUT` | `/v1/autorenewal/tenantconfiguration` | `put-v1-tenant-renewal-configuration` | `—` | `required` | `application/json` | Update the monitoring configuration |
| `revocation-approvals certificaterevocations-approval-rule-create` | `POST` | `/v1/certificates/revocations/approvalrules` | `certificaterevocations_approval_rule_create` | `—` | `required` | `application/json` | Create an approval rule for certificate |
| `revocation-approvals certificaterevocations-approval-rule-delete` | `DELETE` | `/v1/certificates/revocations/approvalrules/{id}` | `certificaterevocations_approval_rule_delete` | `path:id` | `—` | `application/json` | Delete certificate revocation workflow approval ru |
| `revocation-approvals certificaterevocations-approval-rule-get-by-id` | `GET` | `/v1/certificates/revocations/approvalrules/{id}` | `certificaterevocations_approval_rule_getById` | `path:id` | `—` | `application/json` | Retrieve certificate revocation approval rule by |
| `revocation-approvals certificaterevocations-approval-rule-update` | `PUT` | `/v1/certificates/revocations/approvalrules/{id}` | `certificaterevocations_approval_rule_update` | `path:id` | `optional` | `application/json` | Update certificate revocation workflow approval ru |
| `revocation-approvals certificaterevocations-approval-rules-get-all` | `GET` | `/v1/certificates/revocations/approvalrules` | `certificaterevocations_approval_rules_getAll` | `—` | `—` | `application/json` | Get all certificate revocation approval rules |
| `sub-ca-providers subcaprovider-get-all` | `GET` | `/v1/distributedissuers/subcaproviders` | `subcaprovider_getAll` | `—` | `—` | `application/json` | Get the details of all Sub |
| `sub-ca-providers subcaproviders-create` | `POST` | `/v1/distributedissuers/subcaproviders` | `subcaproviders_create` | `—` | `required` | `application/json` | Create a new Sub CA provider |
| `sub-ca-providers subcaproviders-delete` | `DELETE` | `/v1/distributedissuers/subcaproviders/{id}` | `subcaproviders_delete` | `path:id` | `—` | `application/json` | Remove a Sub CA provider |
| `sub-ca-providers subcaproviders-get-by-id` | `GET` | `/v1/distributedissuers/subcaproviders/{id}` | `subcaproviders_getById` | `path:id` | `—` | `application/json` | Get a Sub CA provider details |
| `sub-ca-providers subcaproviders-update` | `PATCH` | `/v1/distributedissuers/subcaproviders/{id}` | `subcaproviders_update` | `path:id` | `required` | `application/json` | Update a Sub CA provider details |
| `tls-server-endpoints certificateinstances-get-all` | `GET` | `/outagedetection/v1/certificateinstances` | `certificateinstances_getAll` | `—` | `—` | `application/json, text/csv` | Retrieve Certificate Instances |
| `tls-server-endpoints certificateinstances-get-by-id` | `GET` | `/outagedetection/v1/certificateinstances/{id}` | `certificateinstances_getById` | `path:id` | `—` | `application/json` | Get a certificate installation details |
| `tls-server-endpoints certificateinstances-search-get-by-expression` | `POST` | `/outagedetection/v1/certificateinstancesearch` | `certificateinstances_search_getByExpression` | `—` | `optional` | `application/json, text/csv` | Retrieve certificate instance data matching search |
| `tls-server-endpoints certificateinstances-validation` | `POST` | `/outagedetection/v1/certificateinstances/validation` | `certificateinstances_validation` | `—` | `required` | `application/json` | Request validation for a set of |
| `vsatellite create-edgeinstances-update` | `POST` | `/v1/edgeinstances/{id}/update` | `create-edgeinstances-update` | `path:id` | `—` | `application/json` | Trigger manual update of Satellite Instance |
| `vsatellite edgeencryptionkeys-get-all` | `GET` | `/v1/edgeencryptionkeys` | `edgeencryptionkeys_getAll` | `—` | `—` | `application/json` | Retrieve Satellite Encryption Keys |
| `vsatellite edgeencryptionkeys-get-by-id` | `GET` | `/v1/edgeencryptionkeys/{id}` | `edgeencryptionkeys_getById` | `path:id` | `—` | `application/json` | Retrieve SatelliteEncryption Key By Id |
| `vsatellite edgeinstances-get-all` | `GET` | `/v1/edgeinstances` | `edgeinstances_getAll` | `—` | `—` | `application/json` | Retrieve Satellite Instances |
| `vsatellite edgeinstances-get-by-id` | `GET` | `/v1/edgeinstances/{id}` | `edgeinstances_getById` | `path:id` | `—` | `application/json` | Retrieve Satellite Instance By Id |
| `vsatellite edgeinstances-update` | `PUT` | `/v1/edgeinstances/{id}` | `edgeinstances_update` | `path:id` | `optional` | `application/json` | Update Satellite Instance |
| `vsatellite edgeworker-delete` | `DELETE` | `/v1/edgeworkers/{id}` | `edgeworker_delete` | `path:id` | `—` | `application/json` | Delete Satellite Worker |
| `vsatellite edgeworkers-create` | `POST` | `/v1/edgeworkers` | `edgeworkers_create` | `—` | `optional` | `application/json` | Create Satellite Worker |
| `vsatellite edgeworkers-get-all` | `GET` | `/v1/edgeworkers` | `edgeworkers_getAll` | `—` | `—` | `application/json` | Retrieve Satellite Workers |
| `vsatellite edgeworkers-pair` | `POST` | `/v1/edgeworkers/{id}/pair` | `edgeworkers_pair` | `path:id` | `optional` | `application/json` | Pair Satellite Worker with Satellite Instance |
| `vsatellite pairingcodes-create` | `POST` | `/v1/pairingcodes/satellite` | `pairingcodes_create` | `—` | `optional` | `application/json` | Create Pairing Code for Satellite Instance |
| `vsatellite recoverycodes-create` | `POST` | `/v1/recoverycodes/satellite` | `recoverycodes_create` | `—` | `optional` | `application/json` | Create Recovery Code for Satellite Instance |
| `vsatellite updatesconfig-get` | `GET` | `/v1/updatesconfig` | `updatesconfig_get` | `—` | `—` | `application/json` | Retrieve Updates configuration |
| `vsatellite updatesconfig-patch` | `PATCH` | `/v1/updatesconfig` | `updatesconfig_patch` | `—` | `required` | `application/json` | Create or Update Configuration |
| `workload-policies policies-create` | `POST` | `/v1/distributedissuers/policies` | `policies_create` | `—` | `optional` | `application/json` | Create a new Workload Issuance policy |
| `workload-policies policies-delete` | `DELETE` | `/v1/distributedissuers/policies/{id}` | `policies_delete` | `path:id` | `—` | `application/json` | Remove a Workload Issuance policy |
| `workload-policies policies-get-all` | `GET` | `/v1/distributedissuers/policies` | `policies_getAll` | `—` | `—` | `application/json` | Get the details of all Workload |
| `workload-policies policies-get-by-id` | `GET` | `/v1/distributedissuers/policies/{id}` | `policies_getById` | `path:id` | `—` | `application/json` | Get a Workload Issuance policy details |
| `workload-policies policies-update` | `PATCH` | `/v1/distributedissuers/policies/{id}` | `policies_update` | `path:id` | `optional` | `application/json` | Update a Workload Issuance policy details |
