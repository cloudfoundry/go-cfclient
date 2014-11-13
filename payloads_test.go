package cfclient

const listOrgsPayload = `{
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
