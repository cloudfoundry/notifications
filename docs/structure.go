package docs

type Endpoint struct {
	Key         string
	Description string
}

type Resource struct {
	Name        string
	Description string
	Endpoints   []Endpoint
}

// Structure defines the top-level structure of the documentation to be generated.
var Structure = []Resource{
	{
		Name:        "Info",
		Description: "EVAN WILL EDIT THIS",
		Endpoints: []Endpoint{
			{
				Key:         "info-get",
				Description: "Retrieve information about the API",
			},
		},
	},
	{
		Name:        "Senders",
		Description: "EVAN WILL EDIT THIS",
		Endpoints: []Endpoint{
			{
				Key:         "sender-create",
				Description: "Create a new sender",
			},
			{
				Key:         "sender-list",
				Description: "List all senders",
			},
			{
				Key:         "sender-get",
				Description: "Retrieve a sender",
			},
			{
				Key:         "sender-update",
				Description: "Update a sender",
			},
			{
				Key:         "sender-delete",
				Description: "Delete a sender",
			},
		},
	},
	{
		Name:        "Templates",
		Description: "EVAN WILL EDIT THIS",
		Endpoints: []Endpoint{
			{
				Key:         "template-create",
				Description: "Create a new template",
			},
			{
				Key:         "template-list",
				Description: "Retrieve a list of templates",
			},
			{
				Key:         "template-get",
				Description: "Retrieve a template",
			},
			{
				Key:         "template-update",
				Description: "Update a template",
			},
			{
				Key:         "template-delete",
				Description: "Delete a template",
			},
		},
	},
	{
		Name:        "Campaign Types",
		Description: "EVAN WILL EDIT THIS",
		Endpoints: []Endpoint{
			{
				Key:         "campaign-type-create",
				Description: "Create a new campaign type",
			},
			{
				Key:         "campaign-type-list",
				Description: "Retrieve a list of campaign types",
			},
			{
				Key:         "campaign-type-get",
				Description: "Retrieve a campaign type",
			},
			{
				Key:         "campaign-type-update",
				Description: "Update a campaign type",
			},
			{
				Key:         "campaign-type-delete",
				Description: "Delete a campaign type",
			},
		},
	},
	{
		Name:        "Campaigns",
		Description: "Campaigns are an email to a set of users using a template provided directly or via a campaign type or via the default template.",
		Endpoints: []Endpoint{
			{
				Key:         "campaign-create",
				Description: "Create a new campaign",
			},
			{
				Key:         "campaign-get",
				Description: "Retrieve a campaign",
			},
			{
				Key:         "campaign-status",
				Description: "Retrieve the status of a campaign",
			},
		},
	},
	{
		Name:        "Unsubscribing",
		Description: "EVAN WILL EDIT THIS",
		Endpoints: []Endpoint{
			{
				Key:         "unsubscriber-put-client",
				Description: "Unsubscribe a user (with a client token)",
			},
			{
				Key:         "unsubscriber-put-user",
				Description: "Unsubscribe a user (with a user token)",
			},
			{
				Key:         "unsubscriber-delete-client",
				Description: "Resubscribe a user (with a client token)",
			},
			{
				Key:         "unsubscriber-delete-user",
				Description: "Resubscribe a user (with a user token)",
			},
		},
	},
}
