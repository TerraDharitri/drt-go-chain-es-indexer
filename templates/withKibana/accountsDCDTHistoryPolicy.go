package withKibana

// AccountsDCDTHistoryPolicy will hold the configuration for the accountsdcdthistory index policy
var AccountsDCDTHistoryPolicy = Object{
	"policy": Object{
		"description":   "Open distro policy for the accountsdcdthistory elastic index.",
		"default_state": "hot",
		"states": Array{
			Object{
				"name": "hot",
				"actions": Array{
					Object{
						"rollover": Object{
							"min_size": "85gb",
						},
					},
				},
				"transitions": Array{
					Object{
						"state_name": "warm",
						"conditions": Object{
							"min_size": "85gb",
						},
					},
				},
			},
			Object{
				"name": "warm",
				"actions": Array{
					Object{
						"replica_count": Object{
							"number_of_replicas": 1,
						},
					},
				},
				"transitions": Array{},
			},
		},
		"ism_template": Object{
			"index_patterns": Array{"accountsdcdthistory-*"},
			"priority":       100,
		},
	},
}
