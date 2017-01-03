package cfclient

const listOrgsPayload = `{
"total_results": 6,
"total_pages": 1,
"prev_url": null,
"next_url": "/v2/orgsPage2",
"resources": [
  {
     "metadata": {
        "guid": "a537761f-9d93-4b30-af17-3d73dbca181b",
        "url": "/v2/organizations/a537761f-9d93-4b30-af17-3d73dbca181b",
        "created_at": "2014-09-24T13:54:53+00:00",
        "updated_at": null
     },
     "entity": {
        "name": "demo",
        "billing_enabled": false,
        "quota_definition_guid": "183599e0-d535-4559-8675-7b6ddb5cc42d",
        "status": "active",
        "quota_definition_url": "/v2/quota_definitions/183599e0-d535-4559-8675-7b6ddb5cc42d",
        "spaces_url": "/v2/organizations/a537761f-9d93-4b30-af17-3d73dbca181b/spaces",
        "domains_url": "/v2/organizations/a537761f-9d93-4b30-af17-3d73dbca181b/domains",
        "private_domains_url": "/v2/organizations/a537761f-9d93-4b30-af17-3d73dbca181b/private_domains",
        "users_url": "/v2/organizations/a537761f-9d93-4b30-af17-3d73dbca181b/users",
        "managers_url": "/v2/organizations/a537761f-9d93-4b30-af17-3d73dbca181b/managers",
        "billing_managers_url": "/v2/organizations/a537761f-9d93-4b30-af17-3d73dbca181b/billing_managers",
        "auditors_url": "/v2/organizations/a537761f-9d93-4b30-af17-3d73dbca181b/auditors",
        "app_events_url": "/v2/organizations/a537761f-9d93-4b30-af17-3d73dbca181b/app_events",
        "space_quota_definitions_url": "/v2/organizations/a537761f-9d93-4b30-af17-3d73dbca181b/space_quota_definitions"
     }
  },
  {
     "metadata": {
        "guid": "da0dba14-6064-4f7a-b15a-ff9e677e49b2",
        "url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2",
        "created_at": "2014-09-26T13:36:41+00:00",
        "updated_at": null
     },
     "entity": {
        "name": "test",
        "billing_enabled": false,
        "quota_definition_guid": "183599e0-d535-4559-8675-7b6ddb5cc42d",
        "status": "active",
        "quota_definition_url": "/v2/quota_definitions/183599e0-d535-4559-8675-7b6ddb5cc42d",
        "spaces_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/spaces",
        "domains_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/domains",
        "private_domains_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/private_domains",
        "users_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/users",
        "managers_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/managers",
        "billing_managers_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/billing_managers",
        "auditors_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/auditors",
        "app_events_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/app_events",
        "space_quota_definitions_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/space_quota_definitions"
     }
  }
]
}`

const listOrgsPayloadPage2 = `{
"total_results": 6,
"total_pages": 1,
"prev_url": null,
"next_url": null,
"resources": [
  {
     "metadata": {
        "guid": "a537761f-9d93-4b30-af17-3d73dbca181b",
        "url": "/v2/organizations/a537761f-9d93-4b30-af17-3d73dbca181b",
        "created_at": "2014-09-24T13:54:53+00:00",
        "updated_at": null
     },
     "entity": {
        "name": "demo",
        "billing_enabled": false,
        "quota_definition_guid": "183599e0-d535-4559-8675-7b6ddb5cc42d",
        "status": "active",
        "quota_definition_url": "/v2/quota_definitions/183599e0-d535-4559-8675-7b6ddb5cc42d",
        "spaces_url": "/v2/organizations/a537761f-9d93-4b30-af17-3d73dbca181b/spaces",
        "domains_url": "/v2/organizations/a537761f-9d93-4b30-af17-3d73dbca181b/domains",
        "private_domains_url": "/v2/organizations/a537761f-9d93-4b30-af17-3d73dbca181b/private_domains",
        "users_url": "/v2/organizations/a537761f-9d93-4b30-af17-3d73dbca181b/users",
        "managers_url": "/v2/organizations/a537761f-9d93-4b30-af17-3d73dbca181b/managers",
        "billing_managers_url": "/v2/organizations/a537761f-9d93-4b30-af17-3d73dbca181b/billing_managers",
        "auditors_url": "/v2/organizations/a537761f-9d93-4b30-af17-3d73dbca181b/auditors",
        "app_events_url": "/v2/organizations/a537761f-9d93-4b30-af17-3d73dbca181b/app_events",
        "space_quota_definitions_url": "/v2/organizations/a537761f-9d93-4b30-af17-3d73dbca181b/space_quota_definitions"
     }
  },
  {
     "metadata": {
        "guid": "da0dba14-6064-4f7a-b15a-ff9e677e49b2",
        "url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2",
        "created_at": "2014-09-26T13:36:41+00:00",
        "updated_at": null
     },
     "entity": {
        "name": "test",
        "billing_enabled": false,
        "quota_definition_guid": "183599e0-d535-4559-8675-7b6ddb5cc42d",
        "status": "active",
        "quota_definition_url": "/v2/quota_definitions/183599e0-d535-4559-8675-7b6ddb5cc42d",
        "spaces_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/spaces",
        "domains_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/domains",
        "private_domains_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/private_domains",
        "users_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/users",
        "managers_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/managers",
        "billing_managers_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/billing_managers",
        "auditors_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/auditors",
        "app_events_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/app_events",
        "space_quota_definitions_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/space_quota_definitions"
     }
  }
]
}`

const orgSpacesPayload = `{
   "total_results": 1,
   "total_pages": 1,
   "prev_url": null,
   "next_url": null,
   "resources": [
      {
         "metadata": {
            "guid": "b8aff561-175d-45e8-b1e7-67e2aedb03b6",
            "url": "/v2/spaces/b8aff561-175d-45e8-b1e7-67e2aedb03b6",
            "created_at": "2014-11-12T17:56:22+00:00",
            "updated_at": null
         },
         "entity": {
            "name": "test",
            "organization_guid": "0c69f181-2d31-4326-ac33-be2b114a5f99",
            "space_quota_definition_guid": null,
            "organization_url": "/v2/organizations/0c69f181-2d31-4326-ac33-be2b114a5f99",
            "developers_url": "/v2/spaces/b8aff561-175d-45e8-b1e7-67e2aedb03b6/developers",
            "managers_url": "/v2/spaces/b8aff561-175d-45e8-b1e7-67e2aedb03b6/managers",
            "auditors_url": "/v2/spaces/b8aff561-175d-45e8-b1e7-67e2aedb03b6/auditors",
            "apps_url": "/v2/spaces/b8aff561-175d-45e8-b1e7-67e2aedb03b6/apps",
            "routes_url": "/v2/spaces/b8aff561-175d-45e8-b1e7-67e2aedb03b6/routes",
            "domains_url": "/v2/spaces/b8aff561-175d-45e8-b1e7-67e2aedb03b6/domains",
            "service_instances_url": "/v2/spaces/b8aff561-175d-45e8-b1e7-67e2aedb03b6/service_instances",
            "app_events_url": "/v2/spaces/b8aff561-175d-45e8-b1e7-67e2aedb03b6/app_events",
            "events_url": "/v2/spaces/b8aff561-175d-45e8-b1e7-67e2aedb03b6/events",
            "security_groups_url": "/v2/spaces/b8aff561-175d-45e8-b1e7-67e2aedb03b6/security_groups"
         }
      }
   ]
}`

const emptyResources = `{
   "total_results": 0,
   "total_pages": 1,
   "prev_url": null,
   "next_url": null,
   "resources": []
}`

const listSpacesPayload = `{
   "total_results": 8,
   "total_pages": 2,
   "prev_url": null,
   "next_url": "/v2/spacesPage2",
   "resources": [
      {
         "metadata": {
            "guid": "8efd7c5c-d83c-4786-b399-b7bd548839e1",
            "url": "/v2/spaces/8efd7c5c-d83c-4786-b399-b7bd548839e1",
            "created_at": "2014-09-24T13:54:54+00:00",
            "updated_at": null
         },
         "entity": {
            "name": "dev",
            "organization_guid": "a537761f-9d93-4b30-af17-3d73dbca181b",
            "space_quota_definition_guid": null,
            "organization_url": "/v2/organizations/a537761f-9d93-4b30-af17-3d73dbca181b",
            "developers_url": "/v2/spaces/8efd7c5c-d83c-4786-b399-b7bd548839e1/developers",
            "managers_url": "/v2/spaces/8efd7c5c-d83c-4786-b399-b7bd548839e1/managers",
            "auditors_url": "/v2/spaces/8efd7c5c-d83c-4786-b399-b7bd548839e1/auditors",
            "apps_url": "/v2/spaces/8efd7c5c-d83c-4786-b399-b7bd548839e1/apps",
            "routes_url": "/v2/spaces/8efd7c5c-d83c-4786-b399-b7bd548839e1/routes",
            "domains_url": "/v2/spaces/8efd7c5c-d83c-4786-b399-b7bd548839e1/domains",
            "service_instances_url": "/v2/spaces/8efd7c5c-d83c-4786-b399-b7bd548839e1/service_instances",
            "app_events_url": "/v2/spaces/8efd7c5c-d83c-4786-b399-b7bd548839e1/app_events",
            "events_url": "/v2/spaces/8efd7c5c-d83c-4786-b399-b7bd548839e1/events",
            "security_groups_url": "/v2/spaces/8efd7c5c-d83c-4786-b399-b7bd548839e1/security_groups"
         }
      },
      {
         "metadata": {
            "guid": "657b5923-7de0-486a-9928-b4d78ee24931",
            "url": "/v2/spaces/657b5923-7de0-486a-9928-b4d78ee24931",
            "created_at": "2014-09-26T13:37:31+00:00",
            "updated_at": null
         },
         "entity": {
            "name": "demo",
            "organization_guid": "da0dba14-6064-4f7a-b15a-ff9e677e49b2",
            "space_quota_definition_guid": null,
            "organization_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2",
            "developers_url": "/v2/spaces/657b5923-7de0-486a-9928-b4d78ee24931/developers",
            "managers_url": "/v2/spaces/657b5923-7de0-486a-9928-b4d78ee24931/managers",
            "auditors_url": "/v2/spaces/657b5923-7de0-486a-9928-b4d78ee24931/auditors",
            "apps_url": "/v2/spaces/657b5923-7de0-486a-9928-b4d78ee24931/apps",
            "routes_url": "/v2/spaces/657b5923-7de0-486a-9928-b4d78ee24931/routes",
            "domains_url": "/v2/spaces/657b5923-7de0-486a-9928-b4d78ee24931/domains",
            "service_instances_url": "/v2/spaces/657b5923-7de0-486a-9928-b4d78ee24931/service_instances",
            "app_events_url": "/v2/spaces/657b5923-7de0-486a-9928-b4d78ee24931/app_events",
            "events_url": "/v2/spaces/657b5923-7de0-486a-9928-b4d78ee24931/events",
            "security_groups_url": "/v2/spaces/657b5923-7de0-486a-9928-b4d78ee24931/security_groups"
         }
      }
   ]
}`

const listSpacesPayloadPage2 = `{
   "total_results": 8,
   "total_pages": 2,
   "prev_url": null,
   "next_url": null,
   "resources": [
      {
         "metadata": {
            "guid": "9ffd7c5c-d83c-4786-b399-b7bd54883977",
            "url": "/v2/spaces/9ffd7c5c-d83c-4786-b399-b7bd54883977",
            "created_at": "2014-09-24T13:54:54+00:00",
            "updated_at": null
         },
         "entity": {
            "name": "test",
            "organization_guid": "a537761f-9d93-4b30-af17-3d73dbca181b",
            "space_quota_definition_guid": null,
            "organization_url": "/v2/organizations/b737761f-9d93-4b30-af17-3d73dbca18aa",
            "developers_url": "/v2/spaces/9ffd7c5c-d83c-4786-b399-b7bd54883977/developers",
            "managers_url": "/v2/spaces/9ffd7c5c-d83c-4786-b399-b7bd54883977/managers",
            "auditors_url": "/v2/spaces/9ffd7c5c-d83c-4786-b399-b7bd54883977/auditors",
            "apps_url": "/v2/spaces/9ffd7c5c-d83c-4786-b399-b7bd54883977/apps",
            "routes_url": "/v2/spaces/9ffd7c5c-d83c-4786-b399-b7bd54883977/routes",
            "domains_url": "/v2/spaces/9ffd7c5c-d83c-4786-b399-b7bd54883977/domains",
            "service_instances_url": "/v2/spaces/9ffd7c5c-d83c-4786-b399-b7bd54883977/service_inst2ances",
            "app_events_url": "/v2/spaces/9ffd7c5c-d83c-4786-b399-b7bd54883977/app_events",
            "events_url": "/v2/spaces/9ffd7c5c-d83c-4786-b399-b7bd54883977/events",
            "security_groups_url": "/v2/spaces/8efd7c5c-d83c-4786-b399-b7bd548839e1/security_groups"
         }
      },
      {
         "metadata": {
            "guid": "329b5923-7de0-486a-9928-b4d78ee24982",
            "url": "/v2/spaces/329b5923-7de0-486a-9928-b4d78ee24982",
            "created_at": "2014-09-26T13:37:31+00:00",
            "updated_at": null
         },
         "entity": {
            "name": "prod",
            "organization_guid": "da0dba14-6064-4f7a-b15a-ff9e677e49b2",
            "space_quota_definition_guid": null,
            "organization_url": "/v2/organizations/ad0dba14-6064-4f7a-b15a-ff9e677e492b",
            "developers_url": "/v2/spaces/329b5923-7de0-486a-9928-b4d78ee24982/developers",
            "managers_url": "/v2/spaces/329b5923-7de0-486a-9928-b4d78ee24982/managers",
            "auditors_url": "/v2/spaces/329b5923-7de0-486a-9928-b4d78ee24982/auditors",
            "apps_url": "/v2/spaces/329b5923-7de0-486a-9928-b4d78ee24982/apps",
            "routes_url": "/v2/spaces/329b5923-7de0-486a-9928-b4d78ee24982/routes",
            "domains_url": "/v2/spaces/329b5923-7de0-486a-9928-b4d78ee24982/domains",
            "service_instances_url": "/v2/spaces/329b5923-7de0-486a-9928-b4d78ee24982/service_instances",
            "app_events_url": "/v2/spaces/329b5923-7de0-486a-9928-b4d78ee24982/app_events",
            "events_url": "/v2/spaces/329b5923-7de0-486a-9928-b4d78ee24982/events",
            "security_groups_url": "/v2/spaces/329b5923-7de0-486a-9928-b4d78ee24982/security_groups"
         }
      }
   ]
}`

const listSecGroupsPayload = `{
   "total_results": 28,
   "total_pages": 1,
   "prev_url": null,
   "next_url": "/v2/security_groupsPage2",
   "resources": [
      {
         "metadata": {
            "guid": "af15c29a-6bde-4a9b-8cdf-43aa0d4b7e3c",
            "url": "/v2/security_groups/af15c29a-6bde-4a9b-8cdf-43aa0d4b7e3c",
            "created_at": "2015-12-04T11:15:55Z",
            "updated_at": null
         },
         "entity": {
            "name": "secgroup-test",
            "rules": [
               {
                  "destination": "1.1.1.1",
                  "ports": "443,4443",
                  "protocol": "tcp"
               },
               {
                  "destination": "1.2.3.4",
                  "ports": "1111",
                  "protocol": "udp"
               }
            ],
            "running_default": true,
            "staging_default": true,
            "spaces_url": "/v2/security_groups/af15c29a-6bde-4a9b-8cdf-43aa0d4b7e3c/spaces",
            "spaces": []
         }
      }
   ]
}`

const listSecGroupsPayloadPage2 = `{
   "total_results": 28,
   "total_pages": 1,
   "prev_url": null,
   "next_url": null,
   "resources": [
      {
         "metadata": {
            "guid": "f9ad202b-76dd-44ec-b7c2-fd2417a561e8",
            "url": "/v2/security_groups/f9ad202b-76dd-44ec-b7c2-fd2417a561e8",
            "created_at": "2015-12-04T11:15:55Z",
            "updated_at": null
         },
         "entity": {
            "name": "secgroup-test2",
            "rules": [
               {
                  "destination": "2.2.2.2",
                  "ports": "2222",
                  "protocol": "udp"
               },
               {
                  "destination": "4.3.2.1",
                  "ports": "443,4443",
                  "protocol": "tcp"
               }
            ],
            "running_default": false,
            "staging_default": false,
            "spaces_url": "/v2/security_groups/f9ad202b-76dd-44ec-b7c2-fd2417a561e8/spaces",
            "spaces": [
               {
                  "metadata": {
                     "guid": "e0a0d1bf-ad74-4b3c-8f4a-0c33859a54e4",
                     "url": "/v2/spaces/e0a0d1bf-ad74-4b3c-8f4a-0c33859a54e4",
                     "created_at": "2014-10-27T10:49:37Z",
                     "updated_at": "2015-01-21T15:30:52Z"
                  },
                  "entity": {
                     "name": "space-test",
                     "organization_guid": "82338ba1-bc08-4576-aad1-9a5b4693b386",
                     "space_quota_definition_guid": null,
                     "allow_ssh": true,
                     "organization_url": "/v2/organizations/82338ba1-bc08-4576-aad1-9a5b4693b386",
                     "developers_url": "/v2/spaces/e0a0d1bf-ad74-4b3c-8f4a-0c33859a54e4/developers",
                     "managers_url": "/v2/spaces/e0a0d1bf-ad74-4b3c-8f4a-0c33859a54e4/managers",
                     "auditors_url": "/v2/spaces/e0a0d1bf-ad74-4b3c-8f4a-0c33859a54e4/auditors",
                     "apps_url": "/v2/spaces/e0a0d1bf-ad74-4b3c-8f4a-0c33859a54e4/apps",
                     "routes_url": "/v2/spaces/e0a0d1bf-ad74-4b3c-8f4a-0c33859a54e4/routes",
                     "domains_url": "/v2/spaces/e0a0d1bf-ad74-4b3c-8f4a-0c33859a54e4/domains",
                     "service_instances_url": "/v2/spaces/e0a0d1bf-ad74-4b3c-8f4a-0c33859a54e4/service_instances",
                     "app_events_url": "/v2/spaces/e0a0d1bf-ad74-4b3c-8f4a-0c33859a54e4/app_events",
                     "events_url": "/v2/spaces/e0a0d1bf-ad74-4b3c-8f4a-0c33859a54e4/events",
                     "security_groups_url": "/v2/spaces/e0a0d1bf-ad74-4b3c-8f4a-0c33859a54e4/security_groups"
                  }
               },
               {
                  "metadata": {
                     "guid": "a2a0d1bf-ad74-4b3c-8f4a-0c33859a5333",
                     "url": "/v2/spaces/a2a0d1bf-ad74-4b3c-8f4a-0c33859a5333",
                     "created_at": "2014-10-27T10:49:37Z",
                     "updated_at": "2015-01-21T15:30:52Z"
                  },
                  "entity": {
                     "name": "space-test2",
                     "organization_guid": "82338ba1-bc08-4576-aad1-9a5b4693b386",
                     "space_quota_definition_guid": null,
                     "allow_ssh": true,
                     "organization_url": "/v2/organizations/82338ba1-bc08-4576-aad1-9a5b4693b386",
                     "developers_url": "/v2/spaces/a2a0d1bf-ad74-4b3c-8f4a-0c33859a5333/developers",
                     "managers_url": "/v2/spaces/a2a0d1bf-ad74-4b3c-8f4a-0c33859a5333/managers",
                     "auditors_url": "/v2/spaces/a2a0d1bf-ad74-4b3c-8f4a-0c33859a5333/auditors",
                     "apps_url": "/v2/spaces/a2a0d1bf-ad74-4b3c-8f4a-0c33859a5333/apps",
                     "routes_url": "/v2/spaces/a2a0d1bf-ad74-4b3c-8f4a-0c33859a5333/routes",
                     "domains_url": "/v2/spaces/a2a0d1bf-ad74-4b3c-8f4a-0c33859a5333/domains",
                     "service_instances_url": "/v2/spaces/a2a0d1bf-ad74-4b3c-8f4a-0c33859a5333/service_instances",
                     "app_events_url": "/v2/spaces/a2a0d1bf-ad74-4b3c-8f4a-0c33859a5333/app_events",
                     "events_url": "/v2/spaces/a2a0d1bf-ad74-4b3c-8f4a-0c33859a5333/events",
                     "security_groups_url": "/v2/spaces/a2a0d1bf-ad74-4b3c-8f4a-0c33859a5333/security_groups"
                  }
               },
               {
                  "metadata": {
                     "guid": "c7a0d1bf-ad74-4b3c-8f4a-0c33859adsa1",
                     "url": "/v2/spaces/c7a0d1bf-ad74-4b3c-8f4a-0c33859adsa1",
                     "created_at": "2014-10-27T10:49:37Z",
                     "updated_at": "2015-01-21T15:30:52Z"
                  },
                  "entity": {
                     "name": "space-test3",
                     "organization_guid": "82338ba1-bc08-4576-aad1-9a5b4693b386",
                     "space_quota_definition_guid": null,
                     "allow_ssh": true,
                     "organization_url": "/v2/organizations/82338ba1-bc08-4576-aad1-9a5b4693b386",
                     "developers_url": "/v2/spaces/c7a0d1bf-ad74-4b3c-8f4a-0c33859adsa1/developers",
                     "managers_url": "/v2/spaces/c7a0d1bf-ad74-4b3c-8f4a-0c33859adsa1/managers",
                     "auditors_url": "/v2/spaces/c7a0d1bf-ad74-4b3c-8f4a-0c33859adsa1/auditors",
                     "apps_url": "/v2/spaces/c7a0d1bf-ad74-4b3c-8f4a-0c33859adsa1/apps",
                     "routes_url": "/v2/spaces/c7a0d1bf-ad74-4b3c-8f4a-0c33859adsa1/routes",
                     "domains_url": "/v2/spaces/c7a0d1bf-ad74-4b3c-8f4a-0c33859adsa1/domains",
                     "service_instances_url": "/v2/spaces/c7a0d1bf-ad74-4b3c-8f4a-0c33859adsa1/service_instances",
                     "app_events_url": "/v2/spaces/c7a0d1bf-ad74-4b3c-8f4a-0c33859adsa1/app_events",
                     "events_url": "/v2/spaces/c7a0d1bf-ad74-4b3c-8f4a-0c33859adsa1/events",
                     "security_groups_url": "/v2/spaces/c7a0d1bf-ad74-4b3c-8f4a-0c33859adsa1/security_groups"
                  }
               }
            ]
         }
      }
   ]
}`

const listAppsPayload = `{
   "total_results": 28,
   "total_pages": 1,
   "prev_url": null,
   "next_url": "/v2/appsPage2",
   "resources": [
      {
         "metadata": {
            "guid": "af15c29a-6bde-4a9b-8cdf-43aa0d4b7e3c",
            "url": "/v2/apps/af15c29a-6bde-4a9b-8cdf-43aa0d4b7e3c",
            "created_at": "2014-10-10T21:03:13+00:00",
            "updated_at": "2014-11-10T14:07:31+00:00"
         },
         "entity": {
            "name": "app-test",
            "production": false,
            "space_guid": "8efd7c5c-d83c-4786-b399-b7bd548839e1",
            "stack_guid": "2c531037-68a2-4e2c-a9e0-71f9d0abf0d4",
            "buildpack": "https://github.com/cloudfoundry/buildpack-go.git",
            "detected_buildpack": null,
            "environment_json": {
               "FOOBAR": "QUX"
            },
            "memory": 256,
            "instances": 1,
            "disk_quota": 1024,
            "state": "STARTED",
            "version": "97ef1272-9eb6-4839-9df1-5ed4f55b5c45",
            "command": null,
            "console": false,
            "debug": null,
            "staging_task_id": "5879c8d06a10491a879734162000def8",
            "package_state": "PENDING",
            "health_check_timeout": null,
            "staging_failed_reason": null,
            "docker_image": null,
            "package_updated_at": "2014-11-10T14:08:50+00:00",
            "detected_start_command": "app-launching-service-broker",
            "space_url": "/v2/spaces/8efd7c5c-d83c-4786-b399-b7bd548839e1",
            "stack_url": "/v2/stacks/2c531037-68a2-4e2c-a9e0-71f9d0abf0d4",
            "events_url": "/v2/apps/af15c29a-6bde-4a9b-8cdf-43aa0d4b7e3c/events",
            "service_bindings_url": "/v2/apps/af15c29a-6bde-4a9b-8cdf-43aa0d4b7e3c/service_bindings",
            "routes_url": "/v2/apps/af15c29a-6bde-4a9b-8cdf-43aa0d4b7e3c/routes"
         }
      }
   ]
}`

const listAppsPayloadPage2 = `{
   "total_results": 28,
   "total_pages": 1,
   "prev_url": null,
   "next_url": null,
   "resources": [
      {
         "metadata": {
            "guid": "f9ad202b-76dd-44ec-b7c2-fd2417a561e8",
            "url": "/v2/apps/f9ad202b-76dd-44ec-b7c2-fd2417a561e8",
            "created_at": "2014-10-10T21:03:13+00:00",
            "updated_at": "2014-11-10T14:07:31+00:00"
         },
         "entity": {
            "name": "app-test2",
            "production": false,
            "space_guid": "8efd7c5c-d83c-4786-b399-b7bd548839e1",
            "stack_guid": "2c531037-68a2-4e2c-a9e0-71f9d0abf0d4",
            "buildpack": "https://github.com/cloudfoundry/buildpack-go.git",
            "detected_buildpack": null,
            "environment_json": {
               "FOOBAR": "QUX"
            },
            "memory": 256,
            "instances": 1,
            "disk_quota": 1024,
            "state": "STARTED",
            "version": "97ef1272-9eb6-4839-9df1-5ed4f55b5c45",
            "command": null,
            "console": false,
            "debug": null,
            "staging_task_id": "5879c8d06a10491a879734162000def8",
            "package_state": "PENDING",
            "health_check_timeout": null,
            "staging_failed_reason": null,
            "docker_image": null,
            "package_updated_at": "2014-11-10T14:08:50+00:00",
            "detected_start_command": "app-launching-service-broker",
            "space_url": "/v2/spaces/8efd7c5c-d83c-4786-b399-b7bd548839e1",
            "stack_url": "/v2/stacks/2c531037-68a2-4e2c-a9e0-71f9d0abf0d4",
            "events_url": "/v2/apps/af15c29a-6bde-4a9b-8cdf-43aa0d4b7e3c/events",
            "service_bindings_url": "/v2/apps/af15c29a-6bde-4a9b-8cdf-43aa0d4b7e3c/service_bindings",
            "routes_url": "/v2/apps/af15c29a-6bde-4a9b-8cdf-43aa0d4b7e3c/routes"
         }
      }
   ]
}`

const appPayload = `{
   "metadata": {
      "guid": "9902530c-c634-4864-a189-71d763cb12e2",
      "url": "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2",
      "created_at": "2014-11-07T23:11:39+00:00",
      "updated_at": "2014-11-07T23:12:03+00:00"
   },
   "entity": {
      "name": "test-env",
      "production": false,
      "space_guid": "a72fa1e8-c694-47b3-85f2-55f61fd00d73",
      "stack_guid": "2c531037-68a2-4e2c-a9e0-71f9d0abf0d4",
      "buildpack": null,
      "detected_buildpack": "Ruby",
      "environment_json": {},
      "memory": 256,
      "instances": 1,
      "disk_quota": 1024,
      "state": "STARTED",
      "version": "0d2f5607-ab6a-4abd-91fe-222cde1ea0f8",
      "command": null,
      "console": false,
      "debug": null,
      "staging_task_id": "46267d4a98ae4f4390aed29975453d60",
      "package_state": "STAGED",
      "health_check_timeout": null,
      "staging_failed_reason": null,
      "docker_image": null,
      "package_updated_at": "2014-11-07T23:12:58+00:00",
      "detected_start_command": "rackup -p $PORT",
      "space_url": "/v2/spaces/a72fa1e8-c694-47b3-85f2-55f61fd00d73",
      "stack_url": "/v2/stacks/2c531037-68a2-4e2c-a9e0-71f9d0abf0d4",
      "events_url": "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2/events",
      "service_bindings_url": "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2/service_bindings",
      "routes_url": "/v2/apps/9902530c-c634-4864-a189-71d763cb12e2/routes"
   }
}`

const appPayloadWithEnvironment_json = `{
   "metadata": {
   },
   "entity": {
      "environment_json": {"string": "string", "int": 1}
   }
}`

const appInstancePayload = `{
   "0": {
      "state": "RUNNING",
      "since": 1455210430.5104606,
      "debug_ip": null,
      "debug_port": null,
      "console_ip": null,
      "console_port": null
   },
   "1": {
      "state": "RUNNING",
      "since": 1455210430.3912115,
      "debug_ip": null,
      "debug_port": null,
      "console_ip": null,
      "console_port": null
   }
}`

const appInstanceUnhealthyPayload = `{
   "0": {
      "state": "RUNNING",
      "since": 1455210430.5104606,
      "debug_ip": null,
      "debug_port": null,
      "console_ip": null,
      "console_port": null
   },
   "1": {
      "state": "STARTING",
      "since": 1455210430.3912115,
      "debug_ip": null,
      "debug_port": null,
      "console_ip": null,
      "console_port": null
   }
}`

const spacePayload = `{
   "metadata": {
      "guid": "a72fa1e8-c694-47b3-85f2-55f61fd00d73",
      "url": "/v2/spaces/a72fa1e8-c694-47b3-85f2-55f61fd00d73",
      "created_at": "2014-11-03T16:47:24+00:00",
      "updated_at": null
   },
   "entity": {
      "name": "test-space",
      "organization_guid": "da0dba14-6064-4f7a-b15a-ff9e677e49b2",
      "space_quota_definition_guid": null,
      "organization_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2",
      "developers_url": "/v2/spaces/a72fa1e8-c694-47b3-85f2-55f61fd00d73/developers",
      "managers_url": "/v2/spaces/a72fa1e8-c694-47b3-85f2-55f61fd00d73/managers",
      "auditors_url": "/v2/spaces/a72fa1e8-c694-47b3-85f2-55f61fd00d73/auditors",
      "apps_url": "/v2/spaces/a72fa1e8-c694-47b3-85f2-55f61fd00d73/apps",
      "routes_url": "/v2/spaces/a72fa1e8-c694-47b3-85f2-55f61fd00d73/routes",
      "domains_url": "/v2/spaces/a72fa1e8-c694-47b3-85f2-55f61fd00d73/domains",
      "service_instances_url": "/v2/spaces/a72fa1e8-c694-47b3-85f2-55f61fd00d73/service_instances",
      "app_events_url": "/v2/spaces/a72fa1e8-c694-47b3-85f2-55f61fd00d73/app_events",
      "events_url": "/v2/spaces/a72fa1e8-c694-47b3-85f2-55f61fd00d73/events",
      "security_groups_url": "/v2/spaces/a72fa1e8-c694-47b3-85f2-55f61fd00d73/security_groups"
   }
}`

const orgPayload = `{
   "metadata": {
      "guid": "da0dba14-6064-4f7a-b15a-ff9e677e49b2",
      "url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2",
      "created_at": "2014-09-26T13:36:41+00:00",
      "updated_at": null
   },
   "entity": {
      "name": "test-org",
      "billing_enabled": false,
      "quota_definition_guid": "183599e0-d535-4559-8675-7b6ddb5cc42d",
      "status": "active",
      "quota_definition_url": "/v2/quota_definitions/183599e0-d535-4559-8675-7b6ddb5cc42d",
      "spaces_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/spaces",
      "domains_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/domains",
      "private_domains_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/private_domains",
      "users_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/users",
      "managers_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/managers",
      "billing_managers_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/billing_managers",
      "auditors_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/auditors",
      "app_events_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/app_events",
      "space_quota_definitions_url": "/v2/organizations/da0dba14-6064-4f7a-b15a-ff9e677e49b2/space_quota_definitions"
   }
}`

const listServicePayload = `{
   "total_results": 22,
   "total_pages": 1,
   "prev_url": null,
   "next_url": null,
   "resources": [
      {
         "metadata": {
            "guid": "a3d76c01-c08a-4505-b06d-8603265682a3",
            "url": "/v2/services/a3d76c01-c08a-4505-b06d-8603265682a3",
            "created_at": "2014-09-24T14:10:51+00:00",
            "updated_at": "2014-10-08T00:06:30+00:00"
         },
         "entity": {
            "label": "nats",
            "provider": null,
            "url": null,
            "description": "NATS is a lightweight cloud messaging system",
            "long_description": null,
            "version": null,
            "info_url": null,
            "active": true,
            "bindable": true,
            "unique_id": "b9310aba-2fa4-11e4-b626-a6c5e4d22fb7",
            "extra": "",
            "tags": [
               "nats",
               "mbus",
               "pubsub"
            ],
            "requires": [],
            "documentation_url": null,
            "service_broker_guid": "a4bdf03a-f0c4-43f9-9c77-f434da91404f",
            "plan_updateable": false,
            "service_plans_url": "/v2/services/a3d76c01-c08a-4505-b06d-8603265682a3/service_plans"
         }
      },
      {
         "metadata": {
            "guid": "ab9ad9c8-1f51-463a-ae3a-5082e9f04ae6",
            "url": "/v2/services/ab9ad9c8-1f51-463a-ae3a-5082e9f04ae6",
            "created_at": "2014-09-24T14:10:51+00:00",
            "updated_at": "2014-10-08T00:06:30+00:00"
         },
         "entity": {
            "label": "etcd",
            "provider": null,
            "url": null,
            "description": "Etcd key-value storage",
            "long_description": null,
            "version": null,
            "info_url": null,
            "active": true,
            "bindable": true,
            "unique_id": "211411a0-2da1-11e4-852f-a6c5e4d22fb7",
            "extra": "",
            "tags": [
               "etcd",
               "keyvalue",
               "etcd-0.4.6"
            ],
            "requires": [],
            "documentation_url": null,
            "service_broker_guid": "a4bdf03a-f0c4-43f9-9c77-f434da91404f",
            "plan_updateable": false,
            "service_plans_url": "/v2/services/ab9ad9c8-1f51-463a-ae3a-5082e9f04ae6/service_plans"
         }
      }
   ]
}`
const listAppsCreatedEventPayload = `{
   "total_results": 2,
   "total_pages": 2,
   "prev_url": null,
   "next_url": "/v2/events2",
   "resources": [
      {
         "metadata": {
            "guid": "49ab122b-82b9-4623-8a13-24e585e32e66",
            "url": "/v2/events/49ab122b-82b9-4623-8a13-24e585e32e66",
            "created_at": "2016-02-26T13:00:21Z",
            "updated_at": null
         },
         "entity": {
            "type": "audit.app.update",
            "actor": "fbf30c43-436e-40e4-8ace-31970b52ce89",
            "actor_type": "user",
            "actor_name": "team-toad@sap.com",
            "actee": "3ca436ff-67a8-468a-8c7d-27ec68a6cfe5",
            "actee_type": "app",
            "actee_name": "authentication-v1-pre-blue",
            "timestamp": "2016-02-26T13:00:21Z",
            "metadata": {
               "request": {
                  "state": "STOPPED"
               }
            },
            "space_guid": "08582a96-cbef-463c-822e-bda8d4284cc7",
            "organization_guid": "bfdcdf09-a3b8-46f4-ab74-d494efefe5b4"
         }
      }
   ]
 }`
const listAppsCreatedEventPayload2 = `{
   "total_results": 2,
   "total_pages": 2,
   "prev_url": "/v2/events",
   "next_url": null,
   "resources": [
      {
         "metadata": {
            "guid": "8e8f83c7-3fc3-4127-9359-ae391380b971",
            "url": "/v2/events/8e8f83c7-3fc3-4127-9359-ae391380b971",
            "created_at": "2016-02-26T13:00:21Z",
            "updated_at": null
         },
         "entity": {
            "type": "audit.app.update",
            "actor": "fbf30c43-436e-40e4-8ace-31970b52ce89",
            "actor_type": "user",
            "actor_name": "team-toad@sap.com",
            "actee": "3ca436ff-67a8-468a-8c7d-27ec68a6cfe5",
            "actee_type": "app",
            "actee_name": "authentication-v1-pre-blue",
            "timestamp": "2016-02-26T13:00:21Z",
            "metadata": {
               "request": {
                  "health_check_timeout": 180,
                  "buildpack": "nodejs_buildpack",
                  "command": "PRIVATE DATA HIDDEN",
                  "state": "STARTED"
               }
            },
            "space_guid": "08582a96-cbef-463c-822e-bda8d4284cc7",
            "organization_guid": "bfdcdf09-a3b8-46f4-ab74-d494efefe5b4"
         }
      }
    ]
 }`

var serviceInstancePayload = `{
   "metadata": {
      "guid": "8423ca96-90ad-411f-b77a-0907844949fc",
      "url": "/v2/service_instances/8423ca96-90ad-411f-b77a-0907844949fc",
      "created_at": "2016-10-21T18:22:56Z",
      "updated_at": null
   },
   "entity": {
      "name": "fortunes-db",
      "credentials": {},
      "service_plan_guid": "f48419f7-4717-4706-86e4-a24973848a77",
      "space_guid": "21e5fdc7-5131-4743-8447-6373cf336a77",
      "gateway_data": null,
      "dashboard_url": "https://p-mysql.system.example.com/manage/instances/8423ca96-90ad-411f-b77a-0907844949fc",
      "type": "managed_service_instance",
      "last_operation": {
         "type": "create",
         "state": "succeeded",
         "description": "",
         "updated_at": null,
         "created_at": "2016-10-21T18:22:56Z"
      },
      "tags": [],
      "space_url": "/v2/spaces/21e5fdc7-5131-4743-8447-6373cf336a77",
      "service_plan_url": "/v2/service_plans/f48419f7-4717-4706-86e4-a24973848a77",
      "service_bindings_url": "/v2/service_instances/8423ca96-90ad-411f-b77a-0907844949fc/service_bindings",
      "service_keys_url": "/v2/service_instances/8423ca96-90ad-411f-b77a-0907844949fc/service_keys",
      "routes_url": "/v2/service_instances/8423ca96-90ad-411f-b77a-0907844949fc/routes"
   }
}`

const listTasksPayload string = `
{
"pagination": {
"total_results": 2,
"total_pages": 1,
"first": {
"href": "https://api.run.example.com/v3/tasks?page=1&per_page=50"
},
"last": {
"href": "https://api.run.example.com/v3/tasks?page=1&per_page=50"
},
"next": null,
"previous": null
},
"resources": [
{
"guid": "xxxxxxxx-e99c-4d60-xxx-e066eb45f8a7",
"sequence_id": 1,
"name": "xxxxxxxx",
"state": "FAILED",
"memory_in_mb": 1024,
"disk_in_mb": 1024,
"result": {
"failure_reason": "Exited with status 127"
},
"created_at": "2016-12-22T13:24:20Z",
"updated_at": "2016-12-22T13:24:25Z",
"droplet_guid": "xxxxxxxx-6cae-49b0-xxxxx-9265950fc16b",
"links": {
"self": {
"href": "https://api.run.example.com/v3/tasks/xxxxxxxx-e99c-4d60-xxxxx-e066eb45f8a7"
},
"app": {
"href": "https://api.run.example.com/v3/apps/xxxxxxxxx-1b30-4e4d-xxxxx-44dec11e3d5b"
},
"droplet": {
"href": "https://api.run.example.com/v3/droplets/xxxxxxxxx-6cae-490b-xxxxx-9265950fc16b"
}
}
},
{
"guid": "xxxxxxxx-5a25-4110-xxx-b309dc5cb0aa",
"sequence_id": 2,
"name": "yyyyyyyyy",
"state": "FAILED",
"memory_in_mb": 1024,
"disk_in_mb": 1024,
"result": {
"failure_reason": "Exited with status 127"
},
"created_at": "2016-12-22T13:24:36Z",
"updated_at": "2016-12-22T13:24:42Z",
"droplet_guid": "xxxxxxx-6cae-49b0-xxxx-9265950fc16b",
"links": {
"self": {
"href": "https://api.run.example.com/v3/tasks/xxxxxxxxx-5a25-4110-xxxxx-b309dc5cb0aa"
},
"app": {
"href": "https://api.run.example.com/v3/apps/xxxxxxxxx-1b30-4e4d-xxxxx-44dec11e3d5b"
},
"droplet": {
"href": "https://api.run.example.com/v3/droplets/xxxxxxxx-6cae-490b-xxxxx-9265950fc16b"
}
}
}
]
}
`

const createTaskPayload = `
{
  "guid": "d5cc22ec-99a3-4e6a-af91-a44b4ab7b6fa",
  "sequence_id": 1,
  "name": "migrate",
  "command": "rake db:migrate",
  "state": "RUNNING",
  "memory_in_mb": 512,
  "disk_in_mb": 1024,
  "result": {
    "failure_reason": null
  },
  "droplet_guid": "740ebd2b-162b-469a-bd72-3edb96fabd9a",
  "created_at": "2016-05-04T17:00:41Z",
  "updated_at": "2016-05-04T17:00:42Z",
  "links": {
    "self": {
      "href": "https://api.example.org/v3/tasks/d5cc22ec-99a3-4e6a-af91-a44b4ab7b6fa"
    },
    "app": {
      "href": "https://api.example.org/v3/apps/ccc25a0f-c8f4-4b39-9f1b-de9f328d0ee5"
    },
    "droplet": {
      "href": "https://api.example.org/v3/droplets/740ebd2b-162b-469a-bd72-3edb96fabd9a"
    }
  }
}
`
