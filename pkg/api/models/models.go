package models

type CreateDeploymentResponse struct {
	Containers []struct {
		Env []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"env"`
		Image           string `json:"image"`
		ImagePullPolicy string `json:"imagePullPolicy"`
		Name            string `json:"name"`
		Ports           []struct {
			ContainerPort int64  `json:"containerPort"`
			Protocol      string `json:"protocol"`
		} `json:"ports"`
		Resources struct {
			Limits struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
			} `json:"limits"`
			Requests struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
			} `json:"requests"`
		} `json:"resources"`
		TerminationMessagePath   string `json:"terminationMessagePath"`
		TerminationMessagePolicy string `json:"terminationMessagePolicy"`
	} `json:"containers"`
	Name string `json:"name"`
}
type CreateDatabaseResponse struct {
	Bootstrap struct {
		Initdb struct {
			Database      string `json:"database"`
			Encoding      string `json:"encoding"`
			LocaleCType   string `json:"localeCType"`
			LocaleCollate string `json:"localeCollate"`
			Owner         string `json:"owner"`
		} `json:"initdb"`
	} `json:"bootstrap"`
	CreationTimestamp     string `json:"creationTimestamp"`
	Instances             int64  `json:"instances"`
	Name                  string `json:"name"`
	Namespace             string `json:"namespace"`
	PrimaryUpdateStrategy string `json:"primaryUpdateStrategy"`
	Services              []struct {
		CreationTimestamp string `json:"creationTimestamp"`
		Name              string `json:"name"`
		Namespace         string `json:"namespace"`
		Ports             []struct {
			Name       string `json:"name"`
			Port       int64  `json:"port"`
			Protocol   string `json:"protocol"`
			TargetPort int64  `json:"targetPort"`
		} `json:"ports"`
	} `json:"services"`
	Storage struct {
		ResizeInUseVolumes bool   `json:"resizeInUseVolumes"`
		Size               string `json:"size"`
	} `json:"storage"`
}

type ListDatabasesResponse struct {
	Bootstrap struct {
		Initdb struct {
			Database      string `json:"database"`
			Encoding      string `json:"encoding"`
			LocaleCType   string `json:"localeCType"`
			LocaleCollate string `json:"localeCollate"`
			Owner         string `json:"owner"`
		} `json:"initdb"`
	} `json:"bootstrap"`
	CreationTimestamp     string `json:"creationTimestamp"`
	Instances             int64  `json:"instances"`
	Name                  string `json:"name"`
	Namespace             string `json:"namespace"`
	PrimaryUpdateStrategy string `json:"primaryUpdateStrategy"`
	Storage               struct {
		ResizeInUseVolumes bool   `json:"resizeInUseVolumes"`
		Size               string `json:"size"`
	} `json:"storage"`
}

type CreateDeploymentInput struct {
	Containers []DeploymentContainerSpec `json:"containers"`
	Name       string                    `json:"name"`
	Replicas   int64                     `json:"replicas"`
}

type DeploymentContainerSpec struct {
	CPU         string        `json:"cpu"`
	Environment []Environment `json:"environment"`

	Image  string `json:"image"`
	Memory string `json:"memory"`
	Name   string `json:"name"`
}

type Environment struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CreateIngressResponse struct {
	Hosts []string `json:"hosts"`
	Name  string   `json:"name"`
}

type DeleteResponse struct {
	Status string `json:"status"`
}
