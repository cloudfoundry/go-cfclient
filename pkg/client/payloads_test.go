package client

const listAppsPayload = `{
  "pagination": {
    "total_results": 2,
    "total_pages": 2,
    "first": {
      "href": "https://api.example.org/v3/apps?page=1&per_page=2"
    },
    "last": {
      "href": "https://api.example.org/v3/apps?page=2&per_page=2"
    },
    "next": {
      "href": "https://api.example.org/v3/apps?page=2&per_page=2"
    },
    "previous": null
  },
  "resources": [
    {
      "guid": "1cb006ee-fb05-47e1-b541-c34179ddc446",
      "name": "my_app",
      "state": "STARTED",
      "created_at": "2016-03-17T21:41:30Z",
      "updated_at": "2016-03-18T11:32:30Z",
      "lifecycle": {
        "type": "buildpack",
        "data": {
          "buildpacks": ["java_buildpack"],
          "stack": "cflinuxfs2"
        }
      },
      "relationships": {
        "space": {
          "data": {
            "guid": "2f35885d-0c9d-4423-83ad-fd05066f8576"
          }
        }
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446"
        },
        "space": {
          "href": "https://api.example.org/v3/spaces/2f35885d-0c9d-4423-83ad-fd05066f8576"
        },
        "processes": {
          "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/processes"
        },
        "route_mappings": {
          "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/route_mappings"
        },
        "packages": {
          "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/packages"
        },
        "environment_variables": {
          "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/environment_variables"
        },
        "current_droplet": {
          "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/droplets/current"
        },
        "droplets": {
          "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/droplets"
        },
        "tasks": {
          "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/tasks"
        },
        "start": {
          "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/actions/start",
          "method": "POST"
        },
        "stop": {
          "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/actions/stop",
          "method": "POST"
        }
      },
      "metadata": {
        "labels": {},
        "annotations": {}
      }
    }
  ]
}`

const listAppsPayloadPage2 = `{
  "pagination": {
    "total_results": 2,
    "total_pages": 2,
    "first": {
      "href": "https://api.example.org/v3/apps?page=1&per_page=2"
    },
    "last": {
      "href": "https://api.example.org/v3/apps?page=2&per_page=2"
    },
    "next": null,
    "previous": null
  },
  "resources": [
    {
      "guid": "02b4ec9b-94c7-4468-9c23-4e906191a0f8",
      "name": "my_app2",
      "state": "STOPPED",
      "created_at": "1970-01-01T00:00:02Z",
      "updated_at": "2016-06-08T16:41:26Z",
      "lifecycle": {
        "type": "buildpack",
        "data": {
          "buildpacks": ["ruby_buildpack", "staticfile_buildpack"],
          "stack": "cflinuxfs2"
        }
      },
      "relationships": {
        "space": {
          "data": {
            "guid": "2f35885d-0c9d-4423-83ad-fd05066f8576"
          }
        }
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/apps/02b4ec9b-94c7-4468-9c23-4e906191a0f8"
        },
        "space": {
          "href": "https://api.example.org/v3/spaces/2f35885d-0c9d-4423-83ad-fd05066f8576"
        },
        "processes": {
          "href": "https://api.example.org/v3/apps/02b4ec9b-94c7-4468-9c23-4e906191a0f8/processes"
        },
        "route_mappings": {
          "href": "https://api.example.org/v3/apps/02b4ec9b-94c7-4468-9c23-4e906191a0f8/route_mappings"
        },
        "packages": {
          "href": "https://api.example.org/v3/apps/02b4ec9b-94c7-4468-9c23-4e906191a0f8/packages"
        },
        "environment_variables": {
          "href": "https://api.example.org/v3/apps/02b4ec9b-94c7-4468-9c23-4e906191a0f8/environment_variables"
        },
        "current_droplet": {
          "href": "https://api.example.org/v3/apps/02b4ec9b-94c7-4468-9c23-4e906191a0f8/droplets/current"
        },
        "droplets": {
          "href": "https://api.example.org/v3/apps/02b4ec9b-94c7-4468-9c23-4e906191a0f8/droplets"
        },
        "tasks": {
          "href": "https://api.example.org/v3/apps/02b4ec9b-94c7-4468-9c23-4e906191a0f8/tasks"
        },
        "start": {
          "href": "https://api.example.org/v3/apps/02b4ec9b-94c7-4468-9c23-4e906191a0f8/actions/start",
          "method": "POST"
        },
        "stop": {
          "href": "https://api.example.org/v3/apps/02b4ec9b-94c7-4468-9c23-4e906191a0f8/actions/stop",
          "method": "POST"
        }
      },
      "metadata": {
        "labels": {},
        "annotations": {}
      }
    }
  ]
}`

const listOrgsPayload = `{
"total_results": 4,
"total_pages": 2,
"prev_url": null,
"next_url": "/v2/organizations?results-per-page=2&page=2",
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

const listSpacesPayload = `{
  "pagination": {
    "total_results": 2,
    "total_pages": 2,
    "first": {
      "href": "https://api.example.org/v3/spaces?page=1&per_page=1"
    },
    "last": {
      "href": "https://api.example.org/v3/spaces?page=2&per_page=1"
    },
    "next": {
      "href": "https://api.example.org/v3/spaces?page=2&per_page=1"
    },
    "previous": null
  },
  "resources": [
    {
      "guid": "space-guid",
      "created_at": "2017-02-01T01:33:58Z",
      "updated_at": "2017-02-01T01:33:58Z",
      "name": "my-space-1",
      "relationships": {
        "organization": {
          "data": {
            "guid": "org-guid"
          }
        },
        "quota": {
          "data": null
        }
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/spaces/space-guid"
        },
        "features": {
          "href": "https://api.example.org/v3/spaces/space-guid/features"
        },
        "organization": {
          "href": "https://api.example.org/v3/organizations/org-guid"
        },
        "apply_manifest": {
          "href": "https://api.example.org/v3/spaces/space-guid/actions/apply_manifest",
          "method": "POST"
        }
      },
      "metadata": {
        "labels": {},
        "annotations": {}
      }
    }
  ]
}`

const listSpacesPayloadPage2 = `{
  "pagination": {
    "total_results": 2,
    "total_pages": 2,
    "first": {
      "href": "https://api.example.org/v3/spaces?page=1&per_page=1"
    },
    "last": {
      "href": "https://api.example.org/v3/spaces?page=2&per_page=1"
    },
    "next": null,
    "previous": {
      "href": "https://api.example.org/v3/spaces?page=2&per_page=1"
    }
  },
  "resources": [
    {
      "guid": "space-guid-2",
      "created_at": "2017-02-01T01:33:58Z",
      "updated_at": "2017-02-01T01:33:58Z",
      "name": "my-space-2",
      "relationships": {
        "organization": {
          "data": {
            "guid": "org-guid"
          }
        },
        "quota": {
          "data": null
        }
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/spaces/space-guid-2"
        },
        "features": {
          "href": "https://api.example.org/v3/spaces/space-guid-2/features"
        },
        "organization": {
          "href": "https://api.example.org/v3/organizations/org-guid"
        },
        "apply_manifest": {
          "href": "https://api.example.org/v3/spaces/space-guid-2/actions/apply_manifest",
          "method": "POST"
        }
      },
      "metadata": {
        "labels": {},
        "annotations": {}
      }
    }
  ]
}`

const listSpaceRolesBySpaceGUIDPayload = `{
  "pagination": {
    "total_results": 3,
    "total_pages": 2,
    "first": {
      "href": "https://api.example.org/v3/roles?page=1&per_page=2&space_guids=spaceGUID1"
    },
    "last": {
      "href": "https://api.example.org/v3/roles?page=2&per_page=2&space_guids=spaceGUID1"
    },
    "next": {
      "href": "https://api.example.org/v3/rolespage2?page=2&per_page=2&space_guids=spaceGUID1"
    },
    "previous": null
  },
  "resources": [
    {
      "guid": "roleGUID1",
      "created_at": "2019-10-10T17:19:12Z",
      "updated_at": "2019-10-10T17:19:12Z",
      "type": "space_developer",
      "relationships": {
        "user": {
          "data": {
            "guid": "userGUID1"
          }
        },
        "space": {
          "data": {
            "guid": "spaceGUID1"
          }
        },
        "organization": {
          "data": null
        }
       },
       "links": {
          "self": {
            "href": "https://api.example.org/v3/roles/roleGUID1"
          },
          "user": {
            "href": "https://api.example.org/v3/users/userGUID1"
          },
          "space": {
            "href": "https://api.example.org/v3/spaces/spaceGUID1"
          }
       }
    },
    {
      "guid": "roleGUID2",
      "created_at": "2047-11-10T17:19:12Z",
      "updated_at": "2047-11-10T17:19:12Z",
      "type": "space_auditor",
      "relationships": {
        "user": {
          "data": {
            "guid": "userGUID2"
          }
        },
        "space": {
          "data": {
            "guid": "spaceGUID1"
          }
        },
        "organization": {
          "data": null
        }
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/roles/roleGUID2"
        },
        "user": {
          "href": "https://api.example.org/v3/users/userGUID2"
        },
        "space": {
          "href": "https://api.example.org/v3/spaces/spaceGUID1"
        }
      }
    }
  ]
}`

const listSpaceRolesBySpaceGuidPayloadPage2 = `{
  "pagination": {
    "total_results": 3,
    "total_pages": 2,
    "first": {
      "href": "https://api.example.org/v3/roles?page=1&per_page=2&space_guids=spaceGUID1"
    },
    "last": {
      "href": "https://api.example.org/v3/rolespage2?page=2&per_page=2&space_guids=spaceGUID1"
    },
    "next": null,
    "previous": {
      "href": "https://api.example.org/v3/roles?page=1&per_page=2&space_guids=spaceGUID1"
    }
  },
  "resources": [
    {
      "guid": "roleGUID3",
      "created_at": "2019-10-10T17:19:12Z",
      "updated_at": "2019-10-10T17:19:12Z",
      "type": "space_manager",
      "relationships": {
        "user": {
          "data": {
            "guid": "userGUID2"
          }
        },
        "space": {
          "data": {
            "guid": "spaceGUID1"
          }
        },
        "organization": {
          "data": null
        }
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/roles/roleGUID3"
        },
        "user": {
          "href": "https://api.example.org/v3/users/userGUID2"
        },
        "space": {
          "href": "https://api.example.org/v3/spaces/spaceGUID1"
        }
      }
    }
  ]
}`

const listSpaceRoleUsersBySpaceGUIDPayload = `{
  "pagination": {
    "total_results": 3,
    "total_pages": 2,
    "first": {
      "href": "https://api.example.org/v3/roles?page=1&per_page=2&include=user&space_guids=spaceGUID1"
    },
    "last": {
      "href": "https://api.example.org/v3/roles?page=2&per_page=2&include=user&space_guids=spaceGUID1"
    },
    "next": {
      "href": "https://api.example.org/v3/rolespage2?page=2&per_page=2&include=user&space_guids=spaceGUID1"
    },
    "previous": null
  },
  "resources": [
    {
      "guid": "roleGUID1",
      "created_at": "2019-10-10T17:19:12Z",
      "updated_at": "2019-10-10T17:19:12Z",
      "type": "space_supporter",
      "relationships": {
        "user": {
          "data": {
            "guid": "userGUID1"
          }
        },
        "space": {
          "data": {
            "guid": "spaceGUID1"
          }
        },
        "organization": {
          "data": null
        }
       },
       "links": {
          "self": {
            "href": "https://api.example.org/v3/roles/roleGUID1"
          },
          "user": {
            "href": "https://api.example.org/v3/users/userGUID1"
          },
          "space": {
            "href": "https://api.example.org/v3/spaces/spaceGUID1"
          }
       }
    },
    {
      "guid": "roleGUID2",
      "created_at": "2047-11-10T17:19:12Z",
      "updated_at": "2047-11-10T17:19:12Z",
      "type": "space_supporter",
      "relationships": {
        "user": {
          "data": {
            "guid": "userGUID2"
          }
        },
        "space": {
          "data": {
            "guid": "spaceGUID1"
          }
        },
        "organization": {
          "data": null
        }
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/roles/roleGUID2"
        },
        "user": {
          "href": "https://api.example.org/v3/users/userGUID2"
        },
        "space": {
          "href": "https://api.example.org/v3/spaces/spaceGUID1"
        }
      }
    }
  ],
  "included": {
    "users": [
      {
        "guid": "userGUID1",
        "created_at": "2022-05-25T23:57:45Z",
        "updated_at": "2022-05-25T23:57:45Z",
        "username": "user1",
        "presentation_name": "user1",
        "origin": "uaa",
        "metadata": {
            "labels": {},
            "annotations": {}
        },
        "links": {
            "self": {
              "href": "https://api.example.org/v3/users/userGUID1"
            }
        }
      },
      {
        "guid": "userGUID2",
        "created_at": "2022-05-25T23:57:45Z",
        "updated_at": "2022-05-25T23:57:45Z",
        "username": "user2",
        "presentation_name": "user2",
        "origin": "uaa",
        "metadata": {
            "labels": {},
            "annotations": {}
        },
        "links": {
            "self": {
              "href": "https://api.example.org/v3/users/userGUID2"
            }
        }
      }
    ]
  }     
}`

const listSpaceRoleUsersBySpaceGUIDPayloadPage2 = `{
  "pagination": {
    "total_results": 3,
    "total_pages": 2,
    "first": {
      "href": "https://api.example.org/v3/roles?page=1&per_page=2&include=user&space_guids=spaceGUID1"
    },
    "last": {
      "href": "https://api.example.org/v3/rolespage2?page=2&per_page=2&include=user&space_guids=spaceGUID1"
    },
    "next": null,
    "previous": {
      "href": "https://api.example.org/v3/roles?page=1&per_page=2&include=user&space_guids=spaceGUID1"
    }
  },
  "resources": [
    {
      "guid": "roleGUID3",
      "created_at": "2019-10-10T17:19:12Z",
      "updated_at": "2019-10-10T17:19:12Z",
      "type": "space_supporter",
      "relationships": {
        "user": {
          "data": {
            "guid": "userGUID2"
          }
        },
        "space": {
          "data": {
            "guid": "spaceGUID1"
          }
        },
        "organization": {
          "data": null
        }
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/roles/roleGUID3"
        },
        "user": {
          "href": "https://api.example.org/v3/users/userGUID2"
        },
        "space": {
          "href": "https://api.example.org/v3/spaces/spaceGUID1"
        }
      }
    }
  ],
  "included": {
    "users": [
      {
        "guid": "userGUID3",
        "created_at": "2022-05-25T23:57:45Z",
        "updated_at": "2022-05-25T23:57:45Z",
        "username": "user3",
        "presentation_name": "user3",
        "origin": "uaa",
        "metadata": {
            "labels": {},
            "annotations": {}
        },
        "links": {
            "self": {
              "href": "https://api.example.org/v3/users/userGUID3"
            }
        }
      }
    ]
  } 
}`

const listSpaceRolesBySpaceGUIDAndTypePayload = `{
  "pagination": {
    "total_results": 3,
    "total_pages": 2,
    "first": {
      "href": "https://api.example.org/v3/roles?page=1&per_page=2&include=user&space_guids=spaceGUID1&types=space_supporter"
    },
    "last": {
      "href": "https://api.example.org/v3/roles?page=2&per_page=2&include=user&space_guids=spaceGUID1&types=space_supporter"
    },
    "next": {
      "href": "https://api.example.org/v3/rolespage2?page=2&per_page=2&include=user&space_guids=spaceGUID1&types=space_supporter"
    },
    "previous": null
  },
  "resources": [
    {
      "guid": "roleGUID1",
      "created_at": "2019-10-10T17:19:12Z",
      "updated_at": "2019-10-10T17:19:12Z",
      "type": "space_supporter",
      "relationships": {
        "user": {
          "data": {
            "guid": "userGUID1"
          }
        },
        "space": {
          "data": {
            "guid": "spaceGUID1"
          }
        },
        "organization": {
          "data": null
        }
       },
       "links": {
          "self": {
            "href": "https://api.example.org/v3/roles/roleGUID1"
          },
          "user": {
            "href": "https://api.example.org/v3/users/userGUID1"
          },
          "space": {
            "href": "https://api.example.org/v3/spaces/spaceGUID1"
          }
       }
    },
    {
      "guid": "roleGUID2",
      "created_at": "2047-11-10T17:19:12Z",
      "updated_at": "2047-11-10T17:19:12Z",
      "type": "space_supporter",
      "relationships": {
        "user": {
          "data": {
            "guid": "userGUID2"
          }
        },
        "space": {
          "data": {
            "guid": "spaceGUID1"
          }
        },
        "organization": {
          "data": null
        }
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/roles/roleGUID2"
        },
        "user": {
          "href": "https://api.example.org/v3/users/userGUID2"
        },
        "space": {
          "href": "https://api.example.org/v3/spaces/spaceGUID1"
        }
      }
    }
  ],
  "included": {
    "users": [
      {
        "guid": "userGUID1",
        "created_at": "2022-05-25T23:57:45Z",
        "updated_at": "2022-05-25T23:57:45Z",
        "username": "user1",
        "presentation_name": "user1",
        "origin": "uaa",
        "metadata": {
            "labels": {},
            "annotations": {}
        },
        "links": {
            "self": {
              "href": "https://api.example.org/v3/users/userGUID1"
            }
        }
      },
      {
        "guid": "userGUID2",
        "created_at": "2022-05-25T23:57:45Z",
        "updated_at": "2022-05-25T23:57:45Z",
        "username": "user2",
        "presentation_name": "user2",
        "origin": "uaa",
        "metadata": {
            "labels": {},
            "annotations": {}
        },
        "links": {
            "self": {
              "href": "https://api.example.org/v3/users/userGUID2"
            }
        }
      }
    ]
  }     
}`

const listSpaceRolesBySpaceGuidAndTypePayloadPage2 = `{
  "pagination": {
    "total_results": 3,
    "total_pages": 2,
    "first": {
      "href": "https://api.example.org/v3/roles?page=1&per_page=2&include=user&space_guids=spaceGUID1&types=space_supporter"
    },
    "last": {
      "href": "https://api.example.org/v3/rolespage2?page=2&per_page=2&include=user&space_guids=spaceGUID1&types=space_supporter"
    },
    "next": null,
    "previous": {
      "href": "https://api.example.org/v3/roles?page=1&per_page=2&include=user&space_guids=spaceGUID1&types=space_supporter"
    }
  },
  "resources": [
    {
      "guid": "roleGUID3",
      "created_at": "2019-10-10T17:19:12Z",
      "updated_at": "2019-10-10T17:19:12Z",
      "type": "space_supporter",
      "relationships": {
        "user": {
          "data": {
            "guid": "userGUID2"
          }
        },
        "space": {
          "data": {
            "guid": "spaceGUID1"
          }
        },
        "organization": {
          "data": null
        }
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/roles/roleGUID3"
        },
        "user": {
          "href": "https://api.example.org/v3/users/userGUID2"
        },
        "space": {
          "href": "https://api.example.org/v3/spaces/spaceGUID1"
        }
      }
    }
  ],
  "included": {
    "users": [
      {
        "guid": "userGUID3",
        "created_at": "2022-05-25T23:57:45Z",
        "updated_at": "2022-05-25T23:57:45Z",
        "username": "user3",
        "presentation_name": "user3",
        "origin": "uaa",
        "metadata": {
            "labels": {},
            "annotations": {}
        },
        "links": {
            "self": {
              "href": "https://api.example.org/v3/users/userGUID3"
            }
        }
      }
    ]
  } 
}`

const listOrganizationRolesByOrganizationGUIDPayload = `{
  "pagination": {
    "total_results": 3,
    "total_pages": 2,
    "first": {
      "href": "https://api.example.org/v3/roles?page=1&per_page=2&include=user&organization_guids=orgGUID1"
    },
    "last": {
      "href": "https://api.example.org/v3/roles?page=2&per_page=2&include=user&organization_guids=orgGUID1"
    },
    "next": {
      "href": "https://api.example.org/v3/rolespage2?page=2&per_page=2&include=user&organization_guids=orgGUID1"
    },
    "previous": null
  },
  "resources": [
    {
      "guid": "roleGUID1",
      "created_at": "2019-10-10T17:19:12Z",
      "updated_at": "2019-10-10T17:19:12Z",
      "type": "organization_auditor",
      "relationships": {
        "user": {
          "data": {
            "guid": "userGUID1"
          }
        },
        "space": {
          "data": null
        },
        "organization": {
          "data": {
            "guid": "orgGUID1"
          }
        }
       },
       "links": {
          "self": {
            "href": "https://api.example.org/v3/roles/roleGUID1"
          },
          "user": {
            "href": "https://api.example.org/v3/users/userGUID1"
          },
          "org": {
            "href": "https://api.example.org/v3/organization/orgGUID1"
          }
       }
    },
    {
      "guid": "roleGUID2",
      "created_at": "2047-11-10T17:19:12Z",
      "updated_at": "2047-11-10T17:19:12Z",
      "type": "organization_auditor",
      "relationships": {
        "user": {
          "data": {
            "guid": "userGUID2"
          }
        },
        "space": {
          "data": null
        },
        "organization": {
          "data": {
            "guid": "orgGUID1"
          }
        }
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/roles/roleGUID2"
        },
        "user": {
          "href": "https://api.example.org/v3/users/userGUID2"
        },
        "organization": {
          "href": "https://api.example.org/v3/organization/orgGUID1"
        }
      }
    }
  ],
  "included": {
    "users": [
      {
        "guid": "userGUID1",
        "created_at": "2022-05-25T23:57:45Z",
        "updated_at": "2022-05-25T23:57:45Z",
        "username": "user1",
        "presentation_name": "user1",
        "origin": "uaa",
        "metadata": {
            "labels": {},
            "annotations": {}
        },
        "links": {
            "self": {
              "href": "https://api.example.org/v3/users/userGUID1"
            }
        }
      },
      {
        "guid": "userGUID2",
        "created_at": "2022-05-25T23:57:45Z",
        "updated_at": "2022-05-25T23:57:45Z",
        "username": "user2",
        "presentation_name": "user2",
        "origin": "uaa",
        "metadata": {
            "labels": {},
            "annotations": {}
        },
        "links": {
            "self": {
              "href": "https://api.example.org/v3/users/userGUID2"
            }
        }
      }
    ]
  }     
}`

const listOrganizationRolesByOrganizationGuidPayloadPage2 = `{
  "pagination": {
    "total_results": 3,
    "total_pages": 2,
    "first": {
      "href": "https://api.example.org/v3/roles?page=1&per_page=2&include=user&organiziation_guids=orgGUID1"
    },
    "last": {
      "href": "https://api.example.org/v3/rolespage2?page=2&per_page=2&include=user&organiziation_guids=orgGUID1"
    },
    "next": null,
    "previous": {
      "href": "https://api.example.org/v3/roles?page=1&per_page=2&include=user&organiziation_guids=orgGUID1"
    }
  },
  "resources": [
    {
      "guid": "roleGUID3",
      "created_at": "2019-10-10T17:19:12Z",
      "updated_at": "2019-10-10T17:19:12Z",
      "type": "organization_auditor",
      "relationships": {
        "user": {
          "data": {
            "guid": "userGUID2"
          }
        },
        "space": {
          "data": null
        },
        "organization": {
          "data": {
            "guid": "spaceGUID1"
          }
        }
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/roles/roleGUID3"
        },
        "user": {
          "href": "https://api.example.org/v3/users/userGUID2"
        },
        "organization": {
          "href": "https://api.example.org/v3/organization/orgGUID1"
        }
      }
    }
  ],
  "included": {
    "users": [
      {
        "guid": "userGUID3",
        "created_at": "2022-05-25T23:57:45Z",
        "updated_at": "2022-05-25T23:57:45Z",
        "username": "user3",
        "presentation_name": "user3",
        "origin": "uaa",
        "metadata": {
            "labels": {},
            "annotations": {}
        },
        "links": {
            "self": {
              "href": "https://api.example.org/v3/users/userGUID3"
            }
        }
      }
    ]
  } 
}`

const listOrganizationRolesByOrganizationGUIDAndTypePayload = `{
  "pagination": {
    "total_results": 3,
    "total_pages": 2,
    "first": {
      "href": "https://api.example.org/v3/roles?page=1&per_page=2&include=user&organization_guids=orgGUID1&types=organization_auditor"
    },
    "last": {
      "href": "https://api.example.org/v3/roles?page=2&per_page=2&include=user&organization_guids=orgGUID1&types=organization_auditor"
    },
    "next": {
      "href": "https://api.example.org/v3/rolespage2?page=2&per_page=2&include=user&organization_guids=orgGUID1&types=organization_auditor"
    },
    "previous": null
  },
  "resources": [
    {
      "guid": "roleGUID1",
      "created_at": "2019-10-10T17:19:12Z",
      "updated_at": "2019-10-10T17:19:12Z",
      "type": "organization_auditor",
      "relationships": {
        "user": {
          "data": {
            "guid": "userGUID1"
          }
        },
        "space": {
          "data": null
        },
        "organization": {
          "data": {
            "guid": "orgGUID1"
          }
        }
       },
       "links": {
          "self": {
            "href": "https://api.example.org/v3/roles/roleGUID1"
          },
          "user": {
            "href": "https://api.example.org/v3/users/userGUID1"
          },
          "org": {
            "href": "https://api.example.org/v3/organization/orgGUID1"
          }
       }
    },
    {
      "guid": "roleGUID2",
      "created_at": "2047-11-10T17:19:12Z",
      "updated_at": "2047-11-10T17:19:12Z",
      "type": "organization_auditor",
      "relationships": {
        "user": {
          "data": {
            "guid": "userGUID2"
          }
        },
        "space": {
          "data": null
        },
        "organization": {
          "data": {
            "guid": "orgGUID1"
          }
        }
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/roles/roleGUID2"
        },
        "user": {
          "href": "https://api.example.org/v3/users/userGUID2"
        },
        "organization": {
          "href": "https://api.example.org/v3/organization/orgGUID1"
        }
      }
    }
  ],
  "included": {
    "users": [
      {
        "guid": "userGUID1",
        "created_at": "2022-05-25T23:57:45Z",
        "updated_at": "2022-05-25T23:57:45Z",
        "username": "user1",
        "presentation_name": "user1",
        "origin": "uaa",
        "metadata": {
            "labels": {},
            "annotations": {}
        },
        "links": {
            "self": {
              "href": "https://api.example.org/v3/users/userGUID1"
            }
        }
      },
      {
        "guid": "userGUID2",
        "created_at": "2022-05-25T23:57:45Z",
        "updated_at": "2022-05-25T23:57:45Z",
        "username": "user2",
        "presentation_name": "user2",
        "origin": "uaa",
        "metadata": {
            "labels": {},
            "annotations": {}
        },
        "links": {
            "self": {
              "href": "https://api.example.org/v3/users/userGUID2"
            }
        }
      }
    ]
  }     
}`

const listOrganizationRolesByOrganizationGuidAndTypePayloadPage2 = `{
  "pagination": {
    "total_results": 3,
    "total_pages": 2,
    "first": {
      "href": "https://api.example.org/v3/roles?page=1&per_page=2&include=user&organiziation_guids=orgGUID1&types=organization_auditor"
    },
    "last": {
      "href": "https://api.example.org/v3/rolespage2?page=2&per_page=2&include=user&organiziation_guids=orgGUID1&types=organization_auditor"
    },
    "next": null,
    "previous": {
      "href": "https://api.example.org/v3/roles?page=1&per_page=2&include=user&organiziation_guids=orgGUID1&types=organization_auditor"
    }
  },
  "resources": [
    {
      "guid": "roleGUID3",
      "created_at": "2019-10-10T17:19:12Z",
      "updated_at": "2019-10-10T17:19:12Z",
      "type": "organization_auditor",
      "relationships": {
        "user": {
          "data": {
            "guid": "userGUID2"
          }
        },
        "space": {
          "data": null
        },
        "organization": {
          "data": {
            "guid": "spaceGUID1"
          }
        }
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/roles/roleGUID3"
        },
        "user": {
          "href": "https://api.example.org/v3/users/userGUID2"
        },
        "organization": {
          "href": "https://api.example.org/v3/organization/orgGUID1"
        }
      }
    }
  ],
  "included": {
    "users": [
      {
        "guid": "userGUID3",
        "created_at": "2022-05-25T23:57:45Z",
        "updated_at": "2022-05-25T23:57:45Z",
        "username": "user3",
        "presentation_name": "user3",
        "origin": "uaa",
        "metadata": {
            "labels": {},
            "annotations": {}
        },
        "links": {
            "self": {
              "href": "https://api.example.org/v3/users/userGUID3"
            }
        }
      }
    ]
  } 
}`

const listSpaceRolesByUserGuidPayload = `{
  "pagination": {
    "total_results": 2,
    "total_pages": 1,
    "first": {
      "href": "https://api.example.org/v3/roles?page=1&per_page=2&user_guids=userGUID1"
    },
    "last": {
      "href": "https://api.example.org/v3/roles?page=2&per_page=2&user_guids=userGUID1"
    },
    "next": null,
    "previous": null
  },
  "resources": [
    {
      "guid": "roleGUID1",
      "created_at": "2019-10-10T17:19:12Z",
      "updated_at": "2019-10-10T17:19:12Z",
      "type": "space_developer",
      "relationships": {
        "user": {
          "data": {
            "guid": "userGUID1"
          }
        },
        "space": {
          "data": {
            "guid": "spaceGUID1"
          }
        },
        "organization": {
          "data": null
        }
       },
       "links": {
          "self": {
            "href": "https://api.example.org/v3/roles/roleGUID1"
          },
          "user": {
            "href": "https://api.example.org/v3/users/userGUID1"
          },
          "space": {
            "href": "https://api.example.org/v3/spaces/spaceGUID1"
          }
       }
    },
    {
      "guid": "roleGUID4",
      "created_at": "2019-10-10T17:19:12Z",
      "updated_at": "2019-10-10T17:19:12Z",
      "type": "space_manager",
      "relationships": {
        "user": {
          "data": {
            "guid": "userGUID1"
          }
        },
        "space": {
          "data": {
            "guid": "spaceGUID2"
          }
        },
        "organization": {
          "data": null
        }
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/roles/roleGUID4"
        },
        "user": {
          "href": "https://api.example.org/v3/users/userGUID1"
        },
        "space": {
          "href": "https://api.example.org/v3/spaces/spaceGUID2"
        }
      }
    }
  ]
}`

const listSpaceRolesBySpaceAndUserGuidPayload = `{
  "pagination": {
    "total_results": 1,
    "total_pages": 1,
    "first": {
      "href": "https://api.example.org/v3/roles?page=1&per_page=1&space_guids=spaceGUID2&user_guids=userGUID1"
    },
    "last": {
      "href": "https://api.example.org/v3/roles?page=1&per_page=1&space_guids=spaceGUID2&user_guids=userGUID1"
    },
    "next": null,
    "previous": null
  },
  "resources": [
    {
      "guid": "roleGUID4",
      "created_at": "2019-10-10T17:19:12Z",
      "updated_at": "2019-10-10T17:19:12Z",
      "type": "space_manager",
      "relationships": {
        "user": {
          "data": {
            "guid": "userGUID1"
          }
        },
        "space": {
          "data": {
            "guid": "spaceGUID2"
          }
        },
        "organization": {
          "data": null
        }
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/roles/roleGUID4"
        },
        "user": {
          "href": "https://api.example.org/v3/users/userGUID1"
        },
        "space": {
          "href": "https://api.example.org/v3/spaces/spaceGUID2"
        }
      }
    }
  ]
}`

const listSecurityGroupsPayload = `{
  "pagination": {
    "total_results": 2,
    "total_pages": 1,
    "first": {
      "href": "https://api.example.org/v3/security_groups?page=1&per_page=50"
    },
    "last": {
      "href": "https://api.example.org/v3/security_groups?page=1&per_page=50"
    },
    "next": null,
    "previous": null
  },
  "resources": [
    {
      "guid": "guid-1",
      "name": "my-group1",
      "globally_enabled": {
        "running": true,
        "staging": false
      },
      "rules": [
        {
          "protocol": "tcp",
          "destination": "1.2.3.4/10",
          "ports": "443,80,8080"
        },
        {
          "protocol": "icmp",
          "destination": "1.2.3.4/12",
          "type": 8,
          "code": 0,
          "description": "test-desc-1"
        }
      ],
      "relationships": {
        "staging_spaces": {
          "data": [
            { "guid": "space-guid-1" },
            { "guid": "space-guid-2" }
          ]
        },
        "running_spaces": {
          "data": []
        }
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/security_groups/guid-1"
        }
      }
    },
    {
      "guid": "guid-2",
      "name": "my-group2",
      "globally_enabled": {
        "running": false,
        "staging": true
      },
      "rules": [
        {
          "protocol": "tcp",
          "destination": "1.2.3.4/14",
          "ports": "443,80,8080"
        },
        {
          "protocol": "icmp",
          "destination": "1.2.3.4/16",
          "type": 5,
          "code": 0,
          "description": "test-desc-2"
        }
      ],
      "relationships": {
        "staging_spaces": {
          "data": [
            { "guid": "space-guid-3" },
            { "guid": "space-guid-4" }
          ]
        },
        "running_spaces": {
          "data": [
            { "guid": "space-guid-5" },
            { "guid": "space-guid-6" }
          ]
        }
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/security_groups/guid-2"
        }
      }
    }
  ]
}`

const listSecurityGroupsByGuidPayload = `{
  "pagination": {
    "total_results": 1,
    "total_pages": 1,
    "first": {
      "href": "https://api.example.org/v3/security_groups?page=1&per_page=50"
    },
    "last": {
      "href": "https://api.example.org/v3/security_groups?page=1&per_page=50"
    },
    "next": null,
    "previous": null
  },
  "resources": [
    {
      "guid": "guid-1",
      "name": "my-group1",
      "globally_enabled": {
        "running": true,
        "staging": false
      },
      "rules": [
        {
          "protocol": "tcp",
          "destination": "1.2.3.4/10",
          "ports": "443,80,8080"
        },
        {
          "protocol": "icmp",
          "destination": "1.2.3.4/12",
          "type": 8,
          "code": 0,
          "description": "test-desc-1"
        }
      ],
      "relationships": {
        "staging_spaces": {
          "data": [
            { "guid": "space-guid-1" },
            { "guid": "space-guid-2" }
          ]
        },
        "running_spaces": {
          "data": []
        }
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/security_groups/guid-1"
        }
      }
    }
  ]
}`

const genericSecurityGroupPayload = `{
  "guid": "guid-1",
  "name": "my-sec-group",
  "globally_enabled": {
    "running": true,
    "staging": false
  },
  "rules": [
    {
      "protocol": "tcp",
      "destination": "10.10.10.0/24",
      "ports": "443,80,8080"
    },
    {
      "protocol": "icmp",
      "destination": "10.10.11.0/24",
      "type": 8,
      "code": 0,
      "description": "Allow ping requests to private services"
    }
  ],
  "relationships": {
    "staging_spaces": {
      "data": []
    },
    "running_spaces": {
      "data": [
        { "guid": "space-guid-1" },
        { "guid": "space-guid-2" }
      ]
    }
  },
  "links": {
    "self": {
      "href": "https://api.example.org/v3/security_groups/guid-1"
    }
  }
}`

const setAppEnvironmentVariablesPayload = `{
  "var": {
    "RAILS_ENV": "production",
    "DEBUG": "false"
  },
  "links": {
    "self": {
      "href": "https://api.example.org/v3/apps/[guid]/environment_variables"
    },
    "app": {
      "href": "https://api.example.org/v3/apps/[guid]"
    }
  }
}`

const listDomainsPayload = `
  {
    "pagination": {
      "total_results": 1,
      "total_pages": 1,
      "first": {
        "href": "https://api.example.org/v3/domains?page=1&per_page=2"
      },
      "last": {
        "href": "https://api.example.org/v3/domains?page=2&per_page=2"
      },
      "next": null,
      "previous": null
    },
    "resources": [
      {
		"guid": "3a5d3d89-3f89-4f05-8188-8a2b298c79d5",
		"created_at": "2019-03-08T01:06:19Z",
		"updated_at": "2019-03-08T01:06:19Z",
		"name": "test-domain.com",
		"internal": false,
		"metadata": {
		  "labels": { },
		  "annotations": { }
		},
		"relationships": {
		  "organization": {
			"data": { "guid": "3a3f3d89-3f89-4f05-8188-751b298c79d5" }
		  },
		  "shared_organizations": {
			"data": [
			  {"guid": "404f3d89-3f89-6z72-8188-751b298d88d5"},
			  {"guid": "416d3d89-3f89-8h67-2189-123b298d3592"}
			]
		  }
		},
		"links": {
		  "self": {
			"href": "https://api.example.org/v3/domains/3a5d3d89-3f89-4f05-8188-8a2b298c79d5"
		  },
		  "organization": {
			"href": "https://api.example.org/v3/organizations/3a3f3d89-3f89-4f05-8188-751b298c79d5"
		  },
		  "route_reservations": {
			"href": "https://api.example.org/v3/domains/3a5d3d89-3f89-4f05-8188-8a2b298c79d5/route_reservations"
		  },
		  "shared_organizations": {
			"href": "https://api.example.org/v3/domains/3a5d3d89-3f89-4f05-8188-8a2b298c79d5/relationships/shared_organizations"
		  }
		}
	  }
    ]
}`

const listUsersPayload = `{
  "pagination": {
     "total_results": 3,
     "total_pages": 2,
     "first": {
        "href": "https://api.example.org/v3/users?page=1&per_page=2"
     },
     "last": {
        "href": "https://api.example.org/v3/users?page=2&per_page=2"
     },
     "next": {
        "href": "https://api.example.org/v3/userspage2?page=2&per_page=2"
     },
     "previous": null
  },
  "resources": [
     {
        "guid": "16f43d50-43a2-4981-bae8-633e8248a637",
        "created_at": "2022-08-02T21:37:52Z",
        "updated_at": "2022-08-02T21:37:52Z",
        "username": "smoke_tests",
        "presentation_name": "smoke_tests",
        "origin": "uaa",
        "metadata": {
           "labels": {},
           "annotations": {}
        },
        "links": {
           "self": {
              "href": "https://api.example.org/v3/users/16f43d50-43a2-4981-bae8-633e8248a637"
           }
        }
     },
     {
        "guid": "test1",
        "created_at": "2022-08-02T21:40:34Z",
        "updated_at": "2022-08-02T21:40:34Z",
        "username": "test1",
        "presentation_name": "test1",
        "origin": "uaa",
        "metadata": {
           "labels": {},
           "annotations": {}
        },
        "links": {
           "self": {
              "href": "https://api.example.org/v3/users/test1"
           }
        }
     }
  ]
}`

const listUsersPayloadPage2 = `{
  "pagination": {
     "total_results": 3,
     "total_pages": 2,
     "first": {
        "href": "https://api.example.org/v3/users?page=1&per_page=2"
     },
     "last": {
        "href": "https://api.example.org/v3/users?page=2&per_page=2"
     },
     "next": {
        "href": ""
     },
     "previous": {
        "href": "https://api.example.org/v3/users?page=1&per_page=2"
     }
  },
  "resources": [
     {
        "guid": "test2",
        "created_at": "2022-08-02T21:41:59Z",
        "updated_at": "2022-08-02T21:41:59Z",
        "username": "test2",
        "presentation_name": "test2",
        "origin": "uaa",
        "metadata": {
           "labels": {},
           "annotations": {}
        },
        "links": {
           "self": {
              "href": "https://api.example.org/v3/users/test2"
           }
        }
     }
  ]
}`

const listOrganizationsPayload = `{
  "pagination": {
    "total_results": 2,
    "total_pages": 2,
    "first": {
      "href": "https://api.example.org/v3/organizations?page=1&per_page=1"
    },
    "last": {
      "href": "https://api.example.org/v3/organizations?page=2&per_page=1"
    },
    "next": {
      "href": "https://api.example.org/v3/organizations?page=2&per_page=1"
    },
    "previous": null
  },
  "resources": [
    {
      "guid": "org-guid",
      "created_at": "2017-02-01T01:33:58Z",
      "updated_at": "2017-02-01T01:33:58Z",
      "name": "my-org-1",
      "relationships": {
        "quota": {
          "data": {
            "guid": "quota-guid"
          }
        }
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/organizations/org-guid"
        },
        "domains": {
          "href": "https://api.example.org/v3/organizations/org-guid/domains"
        },
        "default_domain": {
          "href": "https://api.example.org/v3/organizations/org-guid/domains/default"
        },
        "quota": {
          "href": "https://api.example.org/v3/organization_quotas/quota-guid"
        }
      },
      "metadata": {
        "labels": {},
        "annotations": {}
      }
    }
  ]
}`

const listOrganizationsPayloadPage2 = `{
  "pagination": {
    "total_results": 2,
    "total_pages": 2,
    "first": {
      "href": "https://api.example.org/v3/organizations?page=1&per_page=1"
    },
    "last": {
      "href": "https://api.example.org/v3/organizations?page=2&per_page=1"
    },
    "next": null,
    "previous": {
      "href": "https://api.example.org/v3/organizations?page=2&per_page=1"
    }
  },
  "resources": [
    {
      "guid": "org-guid-2",
      "created_at": "2017-02-01T01:33:58Z",
      "updated_at": "2017-02-01T01:33:58Z",
      "name": "my-org-2",
      "relationships": {
        "quota": {
          "data": null
        }
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/organizations/org-guid-2"
        },
        "domains": {
          "href": "https://api.example.org/v3/organizations/org-guid-2/domains"
        },
        "default_domain": {
          "href": "https://api.example.org/v3/organizations/org-guid-2/domains/default"
        }
      },
      "metadata": {
        "labels": {},
        "annotations": {}
      }
    }
  ]
}`

const updateOrganizationPayload = `
{
  "guid": "org-guid",
  "created_at": "2017-02-01T01:33:58Z",
  "updated_at": "2017-02-01T01:33:58Z",
  "name": "my-org",
  "suspended": false,
  "relationships": {
    "quota": {
      "data": null
    }
  },
  "links": {
    "self": {
      "href": "https://api.example.org/v3/organizations/org-guid"
    },
    "domains": {
      "href": "https://api.example.org/v3/organizations/org-guid/domains"
    },
    "default_domain": {
      "href": "https://api.example.org/v3/organizations/org-guid/domains/default"
    }
  },
  "metadata": {
    "labels": {
      "ORG_KEY": "org_value"
    },
    "annotations": {}
  }
}`

const getOrganizationPayload = `{
  "guid": "org-guid",
  "created_at": "2017-02-01T01:33:58Z",
  "updated_at": "2017-02-01T01:33:58Z",
  "name": "my-org",
  "relationships": {
    "quota": {
      "data": {
        "guid": "quota-guid"
      }
    }
  },
  "links": {
    "self": {
      "href": "https://api.example.org/v3/organizations/org-guid"
    },
    "domains": {
      "href": "https://api.example.org/v3/organizations/org-guid/domains"
    },
    "default_domain": {
      "href": "https://api.example.org/v3/organizations/org-guid/domains/default"
    }
  },
  "metadata": {
    "labels": {
      "ORG_KEY": "org_value"
    },
    "annotations": {}
  }
}`

const createOrganizationPayload = `{
  "guid": "org-guid",
  "created_at": "2017-02-01T01:33:58Z",
  "updated_at": "2017-02-01T01:33:58Z",
  "name": "my-org",
  "relationships": {
    "quota": {
      "data": {
        "guid": "quota-guid"
      }
    }
  },
  "links": {
    "self": {
      "href": "https://api.example.org/v3/organizations/org-guid"
    },
    "domains": {
      "href": "https://api.example.org/v3/organizations/org-guid/domains"
    },
    "default_domain": {
      "href": "https://api.example.org/v3/organizations/org-guid/domains/default"
    }
  },
  "metadata": {
    "labels": {
      "ORG_KEY": "org_value"
    },
    "annotations": {}
  }
}`

const listSpaceUsersPayload = `{
  "pagination": {
    "total_results": 2,
    "total_pages": 2,
    "first": {
      "href": "https://api.example.org/v3/spaces/space-guid/users?page=1&per_page=1"
    },
    "last": {
      "href": "https://api.example.org/v3/spaces/space-guid/users?page=2&per_page=1"
    },
    "next": {
      "href": "https://api.example.org/v3/spaces/space-guid/users?page=2&per_page=1"
    },
    "previous": null
  },
  "resources": [
    {
      "guid": "10a93b89-3f89-4f05-7238-8a2b123c79l9",
      "created_at": "2019-03-08T01:06:18Z",
      "updated_at": "2019-03-08T01:06:18Z",
      "username": "some-name-1",
      "presentation_name": "some-name-1",
      "origin": "uaa",
      "metadata": {
        "labels": {},
        "annotations":{}
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/users/10a93b89-3f89-4f05-7238-8a2b123c79l9"
        }
      }
    }
  ]
}`

const listSpaceUsersPayloadPage2 = `{
  "pagination": {
    "total_results": 2,
    "total_pages": 2,
    "first": {
      "href": "https://api.example.org/v3/spaces/space-guid/users?page=1&per_page=1"
    },
    "last": {
      "href": "https://api.example.org/v3/spaces/space-guid/users?page=2&per_page=1"
    },
    "next": null,
    "previous": {
      "href": "https://api.example.org/v3/spaces/space-guid/users?page=1&per_page=1"
    }
  },
  "resources": [
    {
      "guid": "9da93b89-3f89-4f05-7238-8a2b123c79l9",
      "created_at": "2019-03-08T01:06:19Z",
      "updated_at": "2019-03-08T01:06:19Z",
      "username": "some-name-2",
      "presentation_name": "some-name-2",
      "origin": "ldap",
      "metadata": {
        "labels": {},
        "annotations":{}
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/users/9da93b89-3f89-4f05-7238-8a2b123c79l9"
        }
      }
    }
  ]
}`

const updateSpacePayload = `
{
  "guid": "space-guid",
  "created_at": "2017-02-01T01:33:58Z",
  "updated_at": "2017-02-01T01:33:58Z",
  "name": "my-space",
  "relationships": {
    "organization": {
      "data": {
        "guid": "org-guid"
      }
    },
    "quota": {
      "data": null
    }
  },
  "links": {
    "self": {
      "href": "https://api.example.org/v3/spaces/space-guid"
    },
    "features": {
      "href": "https://api.example.org/v3/spaces/space-guid/features"
    },
    "organization": {
      "href": "https://api.example.org/v3/organizations/org-guid"
    },
    "apply_manifest": {
      "href": "https://api.example.org/v3/spaces/space-guid/actions/apply_manifest",
      "method": "POST"
    }
  },
  "metadata": {
    "labels": {
      "SPACE_KEY": "space_value"
    },
    "annotations": {}
  }
}`

const getSpacePayload = `{
  "guid": "space-guid",
  "created_at": "2017-02-01T01:33:58Z",
  "updated_at": "2017-02-01T01:33:58Z",
  "name": "my-space",
  "relationships": {
    "organization": {
      "data": {
        "guid": "org-guid"
      }
    },
    "quota": {
      "data": null
    }
  },
  "links": {
    "self": {
      "href": "https://api.example.org/v3/spaces/space-guid"
    },
    "features": {
      "href": "https://api.example.org/v3/spaces/space-guid/features"
    },
    "organization": {
      "href": "https://api.example.org/v3/organizations/org-guid"
    },
    "apply_manifest": {
      "href": "https://api.example.org/v3/spaces/space-guid/actions/apply_manifest",
      "method": "POST"
    }
  },
  "metadata": {
    "labels": {
      "SPACE_KEY": "space_value"
    },
    "annotations": {}
  }
}`

const createRoutePayload = `{
	"guid": "cbad697f-cac1-48f4-9017-ac08f39dfb31",
	"host": "a-hostname",
	"path": "/some_path",
	"url": "a-hostname.a-domain.com/some_path",
	"created_at": "2019-05-10T17:17:48Z",
	"updated_at": "2019-05-10T17:17:48Z",
	"metadata": {
	  "labels": { "key": "value" },
	  "annotations": { "note": "detailed information" }
	},
	"relationships": {
	  "space": {
		"data": {
		  "guid": "885a8cb3-c07b-4856-b448-eeb10bf36236"
		}
	  },
	  "domain": {
		"data": {
		  "guid": "0b5f3633-194c-42d2-9408-972366617e0e"
		}
	  }
	},
	"links": {
	  "self": {
		"href": "https://api.example.org/v3/routes/cbad697f-cac1-48f4-9017-ac08f39dfb31"
	  },
	  "space": {
		"href": "https://api.example.org/v3/spaces/885a8cb3-c07b-4856-b448-eeb10bf36236"
	  },
	  "domain": {
		"href": "https://api.example.org/v3/domains/0b5f3633-194c-42d2-9408-972366617e0e"
	  },
	  "destinations": {
		"href": "https://api.example.org/v3/routes/cbad697f-cac1-48f4-9017-ac08f39dfb31/destinations"
	  }
	}
}`

const createSpaceRolePayload = `{
   "guid": "b9f59ab2-2b09-438e-bebb-30e8704ffb89",
   "created_at": "2022-05-31T20:14:13Z",
   "updated_at": "2022-05-31T20:14:13Z",
   "type": "space_supporter",
   "relationships": {
      "user": {
         "data": {
            "guid": "c4958204-6b65-43ea-832b-e4c57aea6641"
         }
      },
      "space": {
         "data": {
            "guid": "b40a40c8-58b7-49a0-b47d-9d6fe5d72905"
         }
      },
      "organization": {
         "data": null
      }
   },
   "links": {
      "self": {
         "href": "https://api.example.org/v3/roles/b9f59ab2-2b09-438e-bebb-30e8704ffb89"
      },
      "user": {
         "href": "https://api.example.org/v3/users/c4958204-6b65-43ea-832b-e4c57aea6641"
      },
      "space": {
         "href": "https://api.example.org/v3/spaces/b40a40c8-58b7-49a0-b47d-9d6fe5d72905"
      }
   }
}`

const createOrganizationRolePayload = `{
  "guid": "21cbfaeb-bff7-4cfd-a7a9-6c13ec76f246",
  "created_at": "2022-05-31T18:19:42Z",
  "updated_at": "2022-05-31T18:19:42Z",
  "type": "organization_user",
  "relationships": {
    "user": {
      "data": {
        "guid": "ac2e02c9-2c5c-4712-a620-a68449d263c3"
      }
    },
    "organization": {
      "data": {
        "guid": "fa8a8346-0d92-4729-870c-77ee1934f973"
      }
    },
    "space": {
      "data": null
    }
  },
  "links": {
    "self": {
      "href": "https://api.example.org/v3/roles/21cbfaeb-bff7-4cfd-a7a9-6c13ec76f246"
    },
    "user": {
      "href": "https://api.example.org/v3/users/ac2e02c9-2c5c-4712-a620-a68449d263c3"
    },
    "organization": {
      "href": "https://api.example.org/v3/organizations/fa8a8346-0d92-4729-870c-77ee1934f973"
    }
  }
}`

const createSpacePayload = `{
  "guid": "space-guid",
  "created_at": "2017-02-01T01:33:58Z",
  "updated_at": "2017-02-01T01:33:58Z",
  "name": "my-space",
  "relationships": {
    "organization": {
      "data": {
        "guid": "org-guid"
      }
    },
    "quota": {
      "data": null
    }
  },
  "links": {
    "self": {
      "href": "https://api.example.org/v3/spaces/space-guid"
    },
    "features": {
      "href": "https://api.example.org/v3/spaces/space-guid/features"
    },
    "organization": {
      "href": "https://api.example.org/v3/organizations/org-guid"
    },
    "apply_manifest": {
      "href": "https://api.example.org/v3/spaces/space-guid/actions/apply_manifest",
      "method": "POST"
    }
  },
  "metadata": {
    "labels": {
      "SPACE_KEY": "space_value"
    },
    "annotations": {}
  }
}`

const getAppPayload = `{
  "guid": "1cb006ee-fb05-47e1-b541-c34179ddc446",
  "name": "my_app",
  "state": "STOPPED",
  "created_at": "2016-03-17T21:41:30Z",
  "updated_at": "2016-06-08T16:41:26Z",
  "lifecycle": {
    "type": "buildpack",
    "data": {
      "buildpacks": ["java_buildpack"],
      "stack": "cflinuxfs2"
    }
  },
  "relationships": {
    "space": {
      "data": {
        "guid": "2f35885d-0c9d-4423-83ad-fd05066f8576"
      }
    }
  },
  "links": {
    "self": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446"
    },
    "space": {
      "href": "https://api.example.org/v3/spaces/2f35885d-0c9d-4423-83ad-fd05066f8576"
    },
    "processes": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/processes"
    },
    "route_mappings": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/route_mappings"
    },
    "packages": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/packages"
    },
    "environment_variables": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/environment_variables"
    },
    "current_droplet": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/droplets/current"
    },
    "droplets": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/droplets"
    },
    "tasks": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/tasks"
    },
    "start": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/actions/start",
      "method": "POST"
    },
    "stop": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/actions/stop",
      "method": "POST"
    }
  },
  "metadata": {
    "labels": {},
    "annotations": {
      "contacts": "Bill tel(1111111) email(bill@fixme), Bob tel(222222) pager(3333333#555) email(bob@fixme)"
    }
  }
}`

const getAppEnvPayload = `{
  "staging_env_json": {
    "GEM_CACHE": "http://gem-cache.example.org"
  },
  "running_env_json": {
    "HTTP_PROXY": "http://proxy.example.org"
  },
  "environment_variables": {
    "RAILS_ENV": "production"
  },
  "system_env_json": {
    "VCAP_SERVICES": {
      "mysql": [
        {
          "name": "db-for-my-app",
          "label": "mysql",
          "tags": ["relational", "sql"],
          "plan": "xlarge",
          "credentials": {
            "username": "user",
            "password": "top-secret"
           },
          "syslog_drain_url": "https://syslog.example.org/drain",
          "provider": null
        }
      ]
    }
  },
  "application_env_json": {
    "VCAP_APPLICATION": {
      "limits": {
        "fds": 16384
      },
      "application_name": "my_app",
      "application_uris": [ "my_app.example.org" ],
      "name": "my_app",
      "space_name": "my_space",
      "space_id": "2f35885d-0c9d-4423-83ad-fd05066f8576",
      "uris": [ "my_app.example.org" ],
      "users": null
    }
  }
}`

const createAppPayload = `{
  "guid": "app-guid",
  "name": "my-app",
  "state": "STOPPED",
  "created_at": "2016-03-17T21:41:30Z",
  "updated_at": "2016-06-08T16:41:26Z",
  "lifecycle": {
    "type": "buildpack",
    "data": {
      "buildpacks": ["java_buildpack"],
      "stack": "cflinuxfs2"
    }
  },
  "relationships": {
    "space": {
      "data": {
        "guid": "space-guid"
      }
    }
  },
  "links": {
    "self": {
      "href": "https://api.example.org/v3/apps/app-guid"
    },
    "space": {
      "href": "https://api.example.org/v3/spaces/space-guid"
    },
    "processes": {
      "href": "https://api.example.org/v3/apps/app-guid/processes"
    },
    "route_mappings": {
      "href": "https://api.example.org/v3/apps/app-guid/route_mappings"
    },
    "packages": {
      "href": "https://api.example.org/v3/apps/app-guid/packages"
    },
    "environment_variables": {
      "href": "https://api.example.org/v3/apps/app-guid/environment_variables"
    },
    "current_droplet": {
      "href": "https://api.example.org/v3/apps/app-guid/droplets/current"
    },
    "droplets": {
      "href": "https://api.example.org/v3/apps/app-guid/droplets"
    },
    "tasks": {
      "href": "https://api.example.org/v3/apps/app-guid/tasks"
    },
    "start": {
      "href": "https://api.example.org/v3/apps/app-guid/actions/start",
      "method": "POST"
    },
    "stop": {
      "href": "https://api.example.org/v3/apps/app-guid/actions/stop",
      "method": "POST"
    }
  },
  "metadata": {
    "labels": {},
    "annotations": {}
  }
}`

const createBuildPayload = `{
  "guid": "585bc3c1-3743-497d-88b0-403ad6b56d16",
  "created_at": "2016-03-28T23:39:34Z",
  "updated_at": "2016-06-08T16:41:26Z",
  "created_by": {
    "guid": "3cb4e243-bed4-49d5-8739-f8b45abdec1c",
    "name": "bill",
    "email": "bill@example.com"
  },
  "state": "STAGING",
  "error": null,
  "lifecycle": {
    "type": "buildpack",
    "data": {
      "buildpacks": [ "ruby_buildpack" ],
      "stack": "cflinuxfs2"
    }
  },
  "package": {
    "guid": "8e4da443-f255-499c-8b47-b3729b5b7432"
  },
  "droplet": null,
  "metadata": {
    "labels": { },
    "annotations": { }
  },
  "links": {
    "self": {
      "href": "https://api.example.org/v3/builds/585bc3c1-3743-497d-88b0-403ad6b56d16"
    },
    "app": {
      "href": "https://api.example.org/v3/apps/7b34f1cf-7e73-428a-bb5a-8a17a8058396"
    }
  }
}
`

const startAppPayload = `{
  "guid": "1cb006ee-fb05-47e1-b541-c34179ddc446",
  "name": "my_app",
  "state": "STARTED",
  "created_at": "2016-03-17T21:41:30Z",
  "updated_at": "2016-03-18T11:32:30Z",
  "lifecycle": {
    "type": "buildpack",
    "data": {
      "buildpacks": ["java_buildpack"],
      "stack": "cflinuxfs2"
    }
  },
  "relationships": {
    "space": {
      "data": {
        "guid": "2f35885d-0c9d-4423-83ad-fd05066f8576"
      }
    }
  },
  "links": {
    "self": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446"
    },
    "space": {
      "href": "https://api.example.org/v3/spaces/2f35885d-0c9d-4423-83ad-fd05066f8576"
    },
    "processes": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/processes"
    },
    "route_mappings": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/route_mappings"
    },
    "packages": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/packages"
    },
    "environment_variables": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/environment_variables"
    },
    "current_droplet": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/droplets/current"
    },
    "droplets": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/droplets"
    },
    "tasks": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/tasks"
    },
    "start": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/actions/start",
      "method": "POST"
    },
    "stop": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/actions/stop",
      "method": "POST"
    }
  },
  "metadata": {
    "labels": {},
    "annotations": {}
  }
}`

const updateAppPayload = `{
  "guid": "1cb006ee-fb05-47e1-b541-c34179ddc446",
  "name": "my_app",
  "state": "STARTED",
  "created_at": "2016-03-17T21:41:30Z",
  "updated_at": "2016-03-18T11:32:30Z",
  "lifecycle": {
    "type": "buildpack",
    "data": {
      "buildpacks": ["java_buildpack"],
      "stack": "cflinuxfs2"
    }
  },
  "relationships": {
    "space": {
      "data": {
        "guid": "2f35885d-0c9d-4423-83ad-fd05066f8576"
      }
    }
  },
  "links": {
    "self": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446"
    },
    "space": {
      "href": "https://api.example.org/v3/spaces/2f35885d-0c9d-4423-83ad-fd05066f8576"
    },
    "processes": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/processes"
    },
    "route_mappings": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/route_mappings"
    },
    "packages": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/packages"
    },
    "environment_variables": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/environment_variables"
    },
    "current_droplet": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/droplets/current"
    },
    "droplets": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/droplets"
    },
    "tasks": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/tasks"
    },
    "start": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/actions/start",
      "method": "POST"
    },
    "stop": {
      "href": "https://api.example.org/v3/apps/1cb006ee-fb05-47e1-b541-c34179ddc446/actions/stop",
      "method": "POST"
    }
  },
  "metadata": {
    "labels": {
      "environment": "production",
      "internet-facing": "false"
    },
    "annotations": {}
  }
}`

const currentDropletV3AppPayload = `{
  "data": {
    "guid": "9d8e007c-ce52-4ea7-8a57-f2825d2c6b39"
  },
  "links": {
    "self": {
      "href": "https://api.example.org/v3/apps/d4c91047-7b29-4fda-b7f9-04033e5c9c9f/relationships/current_droplet"
    },
    "related": {
      "href": "https://api.example.org/v3/apps/d4c91047-7b29-4fda-b7f9-04033e5c9c9f/droplets/current"
    }
  }
}`

const getCurrentAppDropletPayload = `{
  "guid": "585bc3c1-3743-497d-88b0-403ad6b56d16",
  "state": "STAGED",
  "error": null,
  "lifecycle": {
    "type": "buildpack",
    "data": {}
  },
  "execution_metadata": "",
  "process_types": {
    "rake": "bundle exec rake",
    "web": "bundle exec rackup config.ru -p $PORT"
  },
  "checksum": {
    "type": "sha256",
    "value": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
  },
  "buildpacks": [
    {
      "name": "ruby_buildpack",
      "detect_output": "ruby 1.6.14",
      "version": "1.1.1.",
      "buildpack_name": "ruby"
    }
  ],
  "stack": "cflinuxfs3",
  "image": null,
  "created_at": "2016-03-28T23:39:34Z",
  "updated_at": "2016-03-28T23:39:47Z",
  "relationships": {
    "app": {
      "data": {
        "guid": "7b34f1cf-7e73-428a-bb5a-8a17a8058396"
      }
    }
  },
  "links": {
    "self": {
      "href": "https://api.example.org/v3/droplets/585bc3c1-3743-497d-88b0-403ad6b56d16"
    },
    "package": {
      "href": "https://api.example.org/v3/packages/8222f76a-9e09-4360-b3aa-1ed329945e92"
    },
    "app": {
      "href": "https://api.example.org/v3/apps/7b34f1cf-7e73-428a-bb5a-8a17a8058396"
    },
    "assign_current_droplet": {
      "href": "https://api.example.org/v3/apps/7b34f1cf-7e73-428a-bb5a-8a17a8058396/relationships/current_droplet",
      "method": "PATCH"
      },
    "download": {
      "href": "https://api.example.org/v3/droplets/585bc3c1-3743-497d-88b0-403ad6b56d16/download"
    }
  },
  "metadata": {
    "labels": {},
    "annotations": {}
  }
}`

const listServiceInstancesPayload = `{
  "pagination": {
    "total_results": 1,
    "total_pages": 1,
    "first": {
      "href": "https://api.example.org/v3/service_instances?page=1&per_page=50"
    },
    "last": {
      "href": "https://api.example.org/v3/service_instances?page=1&per_page=50"
    },
    "next": null,
    "previous": null
  },
  "resources": [
    {
      "guid": "85ccdcad-d725-4109-bca4-fd6ba062b5c8",
      "created_at": "2017-11-17T13:54:21Z",
      "updated_at": "2017-11-17T13:54:21Z",
      "name": "my_service_instance",
      "relationships": {
        "space": {
          "data": {
            "guid": "ae0031f9-dd49-461c-a945-df40e77c39cb"
          }
        }
      },
      "metadata": {
        "labels": { },
        "annotations": { }
      },
      "links": {
        "space": {
          "href": "https://api.example.org/v3/spaces/ae0031f9-dd49-461c-a945-df40e77c39cb"
        }
      }
    }
  ]
}`

const listRoutesPayload = `{
	"pagination": {
	  "total_results": 1,
	  "total_pages": 1,
	  "first": {
		"href": "https://api.example.org/v3/routes?page=1&per_page=1"
	  },
	  "last": {
		"href": "https://api.example.org/v3/routes?page=1&per_page=1"
	  },
	  "next": null,
	  "previous": null
	},
	"resources": [
	  {
		"guid": "cbad697f-cac1-48f4-9017-ac08f39dfb31",
		"host": "a-hostname",
		"path": "/some_path",
		"url": "a-hostname.a-domain.com/some_path",
		"created_at": "2019-05-10T17:17:48Z",
		"updated_at": "2019-05-10T17:17:48Z",
		"metadata": {
		  "labels": {},
		  "annotations": {}
		},
		"relationships": {
		  "space": {
			"data": {
			  "guid": "885a8cb3-c07b-4856-b448-eeb10bf36236"
			}
		  },
		  "domain": {
			"data": {
			  "guid": "0b5f3633-194c-42d2-9408-972366617e0e"
			}
		  }
		},
		"links": {
		  "self": {
			"href": "https://api.example.org/v3/routes/cbad697f-cac1-48f4-9017-ac08f39dfb31"
		  },
		  "space": {
			"href": "https://api.example.org/v3/spaces/885a8cb3-c07b-4856-b448-eeb10bf36236"
		  },
		  "domain": {
			"href": "https://api.example.org/v3/domains/0b5f3633-194c-42d2-9408-972366617e0e"
		  },
		  "destinations": {
			"href": "https://api.example.org/v3/routes/cbad697f-cac1-48f4-9017-ac08f39dfb31/destinations"
		  }
		}
	  }
	]
}`

const listPackagesForAppPayloadPage1 = `{
  "pagination": {
    "total_results": 2,
    "total_pages": 2,
    "first": {
      "href": "https://api.example.org/v3/apps/f2efe391-2b5b-4836-8518-ad93fa9ebf69/packages?page=1&per_page=1"
    },
    "last": {
      "href": "https://api.example.org/v3/apps/f2efe391-2b5b-4836-8518-ad93fa9ebf69/packages?page=2&per_page=1"
    },
    "next": {
      "href": "https://api.example.org/v3/apps/f2efe391-2b5b-4836-8518-ad93fa9ebf69/packages?page=2&per_page=1"
    },
    "previous": null
  },
  "resources": [
    {
      "guid": "752edab0-2147-4f58-9c25-cd72ad8c3561",
      "type": "bits",
      "data": {
        "error": null,
        "checksum": {
          "type": "sha256",
          "value": null
        }
      },
      "state": "READY",
      "created_at": "2016-03-17T21:41:09Z",
      "updated_at": "2016-06-08T16:41:26Z",
      "links": {
        "self": {
          "href": "https://api.example.org/v3/packages/752edab0-2147-4f58-9c25-cd72ad8c3561"
        },
        "upload": {
          "href": "https://api.example.org/v3/packages/752edab0-2147-4f58-9c25-cd72ad8c3561/upload",
          "method": "POST"
        },
        "download": {
          "href": "https://api.example.org/v3/packages/752edab0-2147-4f58-9c25-cd72ad8c3561/download",
          "method": "GET"
        },
        "app": {
          "href": "https://api.example.org/v3/apps/f2efe391-2b5b-4836-8518-ad93fa9ebf69"
        }
      },
      "metadata": {
        "labels": {},
        "annotations": {}
      }
    }
  ]
}`

const listPackagesForAppPayloadPage2 = `{
  "pagination": {
    "total_results": 2,
    "total_pages": 2,
    "first": {
      "href": "https://api.example.org/v3/apps/f2efe391-2b5b-4836-8518-ad93fa9ebf69/packages?page=1&per_page=1"
    },
    "last": {
      "href": "https://api.example.org/v3/apps/f2efe391-2b5b-4836-8518-ad93fa9ebf69/packages?page=2&per_page=1"
    },
    "next": null,
    "previous": null
  },
  "resources": [
    {
      "guid": "2345ab-2147-4f58-9c25-cd72ad8c3561",
      "type": "bits",
      "data": {
        "error": null,
        "checksum": {
          "type": "sha256",
          "value": null
        }
      },
      "state": "READY",
      "created_at": "2016-03-17T21:41:09Z",
      "updated_at": "2016-06-08T16:41:26Z",
      "links": {
        "self": {
          "href": "https://api.example.org/v3/packages/2345ab-2147-4f58-9c25-cd72ad8c3561"
        },
        "upload": {
          "href": "https://api.example.org/v3/packages/2345ab-2147-4f58-9c25-cd72ad8c3561/upload",
          "method": "POST"
        },
        "download": {
          "href": "https://api.example.org/v3/packages/2345ab-2147-4f58-9c25-cd72ad8c3561/download",
          "method": "GET"
        },
        "app": {
          "href": "https://api.example.org/v3/apps/f2efe391-2b5b-4836-8518-ad93fa9ebf69"
        }
      },
      "metadata": {
        "labels": {},
        "annotations": {}
      }
    }
  ]
}`

const copyPackagePayload = `{
  "guid": "fec72fc1-e453-4463-a86d-5df426f337a3",
  "type": "docker",
  "data": {
    "image": "http://awesome-sauce.example.org"
  },
  "state": "COPYING",
  "created_at": "2016-03-17T21:41:09Z",
  "updated_at": "2016-06-08T16:41:26Z",
  "links": {
    "self": {
      "href": "https://api.example.org/v3/packages/fec72fc1-e453-4463-a86d-5df426f337a3"
    },
    "app": {
      "href": "https://api.example.org/v3/apps/36208a68-562d-4f51-94ea-28bd8553a271"
    }
  },
  "metadata": {
    "labels": {},
    "annotations": {}
  }
}`

const getDeploymentPayload = `{
  "guid": "59c3d133-2b83-46f3-960e-7765a129aea4",
  "state": "DEPLOYING",
  "status": {
    "value": "ACTIVE",
    "reason": "DEPLOYING",
    "details": {
      "last_successful_healthcheck": "2018-04-25T22:42:10Z"
    }
  },
  "strategy": "rolling",
  "droplet": {
    "guid": "44ccfa61-dbcf-4a0d-82fe-f668e9d2a962"
  },
  "previous_droplet": {
    "guid": "cc6bc315-bd06-49ce-92c2-bc3ad45268c2"
  },
  "new_processes": [
    {
      "guid": "fd5d3e60-f88c-4c37-b1ae-667cfc65a856",
      "type": "web"
    }
  ],
  "revision": {
    "guid": "56126cba-656a-4eba-a81e-7e9951b2df57",
    "version": 1
  },
  "created_at": "2018-04-25T22:42:10Z",
  "updated_at": "2018-04-25T22:42:10Z",
  "metadata": {
    "labels": { },
    "annotations": { }
  },
  "relationships": {
    "app": {
      "data": {
        "guid": "305cea31-5a44-45ca-b51b-e89c7a8ef8b2"
      }
    }
  },
  "links": {
    "self": {
      "href": "https://api.example.org/v3/deployments/59c3d133-2b83-46f3-960e-7765a129aea4"
    },
    "app": {
      "href": "https://api.example.org/v3/apps/305cea31-5a44-45ca-b51b-e89c7a8ef8b2"
    }
  }
}`

const listStacksPayload = `{
  "pagination": {
    "total_results": 2,
    "total_pages": 1,
    "first": {
      "href": "https://api.example.org/v3/stacks?page=1&per_page=2"
    },
    "last": {
      "href": "https://api.example.org/v3/stacks?page=1&per_page=2"
    },
    "next": null,
    "previous": null
  },
  "resources": [
    {
      "guid": "guid-1",
      "created_at": "2018-11-09T22:43:28Z",
      "updated_at": "2018-11-09T22:43:28Z",
      "name": "my-stack-1",
      "description": "This is my first stack!",
      "metadata": {
        "labels": {},
        "annotations": {}
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/stacks/guid-1"
        }
      }
    },
    {
      "guid": "guid-2",
      "created_at": "2018-11-09T22:43:29Z",
      "updated_at": "2018-11-09T22:43:29Z",
      "name": "my-stack-2",
      "description": "This is my second stack!",
      "metadata": {
        "labels": {},
        "annotations": {}
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/stacks/guid-2"
        }
      }
    }
  ]
}`

const listServiceCredentialBindingsPayload = `{
  "pagination": {
    "total_results": 1,
    "total_pages": 1,
    "first": {
      "href": "https://api.example.org/v3/service_instances?page=1&per_page=50"
    },
    "last": {
      "href": "https://api.example.org/v3/service_instances?page=1&per_page=50"
    },
    "next": null,
    "previous": null
  },
  "resources": [
    {
      "guid": "d9634934-8e1f-4c2d-bb33-fa5df019cf9d",
      "created_at": "2022-02-17T17:17:44Z",
      "updated_at": "2022-02-17T17:17:44Z",
      "name": "my_service_key",
      "type": "key",
      "relationships": {
        "service_instance": {
          "data": {
              "guid": "85ccdcad-d725-4109-bca4-fd6ba062b5c8"
          }
        }
      },
      "metadata": {
        "labels": { },
        "annotations": { }
      },
      "links": {
        "self": {
          "href": "https://api.example.org/v3/service_credential_bindings/d9634934-8e1f-4c2d-bb33-fa5df019cf9d"
        },
        "details": {
          "href": "https://api.example.org/v3/service_credential_bindings/d9634934-8e1f-4c2d-bb33-fa5df019cf9d/details"
        },
        "service_instance": {
          "href": "https://api.example.org/v3/service_instances/85ccdcad-d725-4109-bca4-fd6ba062b5c8"
        },
        "parameters": {
          "href": "https://api.example.org/v3/service_credential_bindings/d9634934-8e1f-4c2d-bb33-fa5df019cf9d/parameters"
        }
      }
    }
  ]
}`

const getServiceCredentialBindingsByGUIDPayload = `{
  "guid": "d9634934-8e1f-4c2d-bb33-fa5df019cf9d",
  "created_at": "2022-02-17T17:17:44Z",
  "updated_at": "2022-02-17T17:17:44Z",
  "name": "my_service_key",
  "type": "key",
  "last_operation": {
    "type": "create",
    "state": "succeeded",
    "description": "",
    "created_at": "2022-02-17T17:17:44Z",
    "updated_at": "2022-02-17T17:17:44Z"
  },
  "relationships": {
    "service_instance": {
      "data": {
        "guid": "85ccdcad-d725-4109-bca4-fd6ba062b5c8"
      }
    }
  },
  "metadata": {
    "labels": { },
    "annotations": { }
  },
  "links": {
    "self": {
      "href": "https://api.example.org/v3/service_credential_bindings/d9634934-8e1f-4c2d-bb33-fa5df019cf9d"
    },
    "details": {
      "href": "https://api.example.org/v3/service_credential_bindings/d9634934-8e1f-4c2d-bb33-fa5df019cf9d/details"
    },
    "service_instance": {
      "href": "https://api.example.org/v3/service_instances/85ccdcad-d725-4109-bca4-fd6ba062b5c8"
    },
    "parameters": {
      "href": "https://api.example.org/v3/service_credential_bindings/d9634934-8e1f-4c2d-bb33-fa5df019cf9d/parameters"
    }
  }
}`
