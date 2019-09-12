package lib

/* import (
	"encoding/json"
) */


type SolutionTemplateConfig struct {
	Repository     int `json:"repository"`
	RequiredValues map[string]string `json:"required_values"`
	Namespace      string `json:"namespace"`
	Values         string `json:"values_yaml"`
	ChartPath      string `json:"chart_path"`
	Logo           string `json:"logo"`
	ReleaseName    string `json:"release_name"`
	ChartName		string `json:"chart_name"`
	Version			string	`json:"version"`
}

type SolutionTemplateSpec struct {
	Requirements SolutionTemplateSpecRequirements `json:"requirements"`
}

type SolutionTemplateSpecRequirements struct {
	Node 	SolutionTemplateSpecRequirementsNode `json:"node"`
	App     SolutionTemplateSpecRequirementsApp `json:"app"`
}

type SolutionTemplateSpecRequirementsNode struct {
	Count	int	`json:"count"`
	CPU	int	`json:"CPU"`
}
type SolutionTemplateSpecRequirementsApp struct {
	RBAC bool	`json:"rbac"`
	Namespace string `json:"namespace"`
	ValuesEditRequired []string `json:"valuesEditRequired"`
	Storage	bool	`json:"storage"`
	LoadBalancer	bool	`json:"loadbalancer"`
}

type SolutionTemplateDependencies struct {
	Name string `json:"name"`
	Value string `json:"value"`
	Available bool `json:"available"`
	KeysetRequired	bool	`json:"keysetRequired"`
	Tag 	string	`json:"tag"`
	IsPostBuildCompatible bool	`json:"isPostBuildCompatible"`
	IsManagedIndependently bool	`json:"isManagedIndependently"`
	Dependencies	[]string	`json:"dependencies"`
}

type SolutionTemplate struct {
	Name        string         `json:"name"`
	Solution    string         `json:"solution"`
	Installer   string         `json:"installer"`
	Keyset      *int         `json:"keyset"`
	Mode	*string `json:"mode"`
	Tag		string 	`json:"tag"`
	Version     string         `json:"version"`
	Config      SolutionTemplateConfig `json:"config"`
	Spec		SolutionTemplateSpec `json:"spec"`
	Dependencies	SolutionTemplateDependencies `json:"dependencies"`
}

var jenkins = SolutionTemplate{
	Name: "jenkins",
	Solution: "jenkins-repo-1",
	Installer: "helm",
	Keyset: nil,
	Mode: nil,
	Tag:	"0.19.1",
	Config: SolutionTemplateConfig{
		Namespace: "jenkins",
		ChartName: "jenkins",
		Version: "0.19.1",
		ChartPath: "spc-trusted-charts/jenkins",
		ReleaseName: "jenkins-dry-recipe",
		Logo: "http://trusted-charts-logos.stackpoint.io/jenkins.png",
		Repository: 1,
		Values: "# Default values for jenkins.\n# This is a YAML-formatted file.\n# Declare name/value pairs to be passed into your templates.\n# name: value\n\n## Overrides for generated resource names\n# See templates/_helpers.tpl\n# nameOverride:\n# fullnameOverride:\n\nMaster:\n  Name: jenkins-master\n  Image: \"jenkins/jenkins\"\n  ImageTag: \"lts\"\n  ImagePullPolicy: \"Always\"\n# ImagePullSecret: jenkins\n  Component: \"jenkins-master\"\n  UseSecurity: true\n  HostNetworking: false\n  AdminUser: admin\n  # AdminPassword: <defaults to random>\n  resources:\n    requests:\n      cpu: \"50m\"\n      memory: \"256Mi\"\n    limits:\n      cpu: \"2000m\"\n      memory: \"2048Mi\"\n  # Environment variables that get added to the init container (useful for e.g. http_proxy)\n  # InitContainerEnv:\n  #   - name: http_proxy\n  #     value: \"http://192.168.64.1:3128\"\n  # ContainerEnv:\n  #   - name: http_proxy\n  #     value: \"http://192.168.64.1:3128\"\n  # Set min/max heap here if needed with:\n  # JavaOpts: \"-Xms512m -Xmx512m\"\n  # JenkinsOpts: \"\"\n  # JenkinsUriPrefix: \"/jenkins\"\n  # Enable pod security context (must be `true` if RunAsUser or FsGroup are set)\n  UsePodSecurityContext: true\n  # Set RunAsUser to 1000 to let Jenkins run as non-root user 'jenkins' which exists in 'jenkins/jenkins' docker image.\n  # When setting RunAsUser to a different value than 0 also set FsGroup to the same value:\n  # RunAsUser: <defaults to 0>\n  # FsGroup: <will be omitted in deployment if RunAsUser is 0>\n  ServicePort: 8080\n  # For minikube, set this to NodePort, elsewhere use LoadBalancer\n  # Use ClusterIP if your setup includes ingress controller\n  ServiceType: LoadBalancer\n  # Master Service annotations\n  ServiceAnnotations: {}\n  #   service.beta.kubernetes.io/aws-load-balancer-backend-protocol: https\n  # Used to create Ingress record (should used with ServiceType: ClusterIP)\n  # HostName: jenkins.cluster.local\n  # NodePort: <to set explicitly, choose port between 30000-32767\n  # Enable Kubernetes Liveness and Readiness Probes\n  # ~ 2 minutes to allow Jenkins to restart when upgrading plugins. Set ReadinessTimeout to be shorter than LivenessTimeout.\n  HealthProbes: true\n  HealthProbesLivenessTimeout: 90\n  HealthProbesReadinessTimeout: 60\n  HealthProbeLivenessFailureThreshold: 12\n  SlaveListenerPort: 50000\n  DisabledAgentProtocols:\n    - JNLP-connect\n    - JNLP2-connect\n  CSRF:\n    DefaultCrumbIssuer:\n      Enabled: true\n      ProxyCompatability: true\n  CLI: false\n  # Kubernetes service type for the JNLP slave service\n  # SETTING THIS TO \"LoadBalancer\" IS A HUGE SECURITY RISK: https://github.com/kubernetes/charts/issues/1341\n  SlaveListenerServiceType: ClusterIP\n  SlaveListenerServiceAnnotations: {}\n  LoadBalancerSourceRanges:\n  - 0.0.0.0/0\n  # Optionally assign a known public LB IP\n  # LoadBalancerIP: 1.2.3.4\n  # Optionally configure a JMX port\n  # requires additional JavaOpts, ie\n  # JavaOpts: >\n  #   -Dcom.sun.management.jmxremote.port=4000\n  #   -Dcom.sun.management.jmxremote.authenticate=false\n  #   -Dcom.sun.management.jmxremote.ssl=false\n  # JMXPort: 4000\n  # List of plugins to be install during Jenkins master start\n  InstallPlugins:\n    - kubernetes:1.12.4\n    - workflow-job:2.24\n    - workflow-aggregator:2.5\n    - credentials-binding:1.16\n    - git:3.9.1\n  # Used to approve a list of groovy functions in pipelines used the script-security plugin. Can be viewed under /scriptApproval\n  # ScriptApproval:\n  #   - \"method groovy.json.JsonSlurperClassic parseText java.lang.String\"\n  #   - \"new groovy.json.JsonSlurperClassic\"\n  # List of groovy init scripts to be executed during Jenkins master start\n  InitScripts:\n  #  - |\n  #    print 'adding global pipeline libraries, register properties, bootstrap jobs...'\n  # Kubernetes secret that contains a 'credentials.xml' for Jenkins\n  # CredentialsXmlSecret: jenkins-credentials\n  # Kubernetes secret that contains files to be put in the Jenkins 'secrets' directory,\n  # useful to manage encryption keys used for credentials.xml for instance (such as\n  # master.key and hudson.util.Secret)\n  # SecretsFilesSecret: jenkins-secrets\n  # Jenkins XML job configs to provision\n  # Jobs: |-\n  #   test: |-\n  #     <<xml here>>\n  CustomConfigMap: false\n  # By default, the configMap is only used to set the initial config the first time\n  # that the chart is installed.  Setting `OverwriteConfig` to `true` will overwrite\n  # the jenkins config with the contents of the configMap every time the pod starts.\n  OverwriteConfig: false\n  # Node labels and tolerations for pod assignment\n  # ref: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector\n  # ref: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#taints-and-tolerations-beta-feature\n  NodeSelector: {}\n  Tolerations: {}\n  PodAnnotations: {}\n\n  Ingress:\n    ApiVersion: extensions/v1beta1\n    Annotations: {}\n    # kubernetes.io/ingress.class: nginx\n    # kubernetes.io/tls-acme: \"true\"\n\n    TLS:\n    # - secretName: jenkins.cluster.local\n    #   hosts:\n    #     - jenkins.cluster.local\n\nAgent:\n  Enabled: true\n  Image: jenkins/jnlp-slave\n  ImageTag: 3.10-1\n  CustomJenkinsLabels: []\n# ImagePullSecret: jenkins\n  Component: \"jenkins-slave\"\n  Privileged: false\n  resources:\n    requests:\n      cpu: \"200m\"\n      memory: \"256Mi\"\n    limits:\n      cpu: \"200m\"\n      memory: \"256Mi\"\n  # You may want to change this to true while testing a new image\n  AlwaysPullImage: false\n  # Controls how slave pods are retained after the Jenkins build completes\n  # Possible values: Always, Never, OnFailure\n  PodRetention: Never\n  # You can define the volumes that you want to mount for this container\n  # Allowed types are: ConfigMap, EmptyDir, HostPath, Nfs, Pod, Secret\n  # Configure the attributes as they appear in the corresponding Java class for that type\n  # https://github.com/jenkinsci/kubernetes-plugin/tree/master/src/main/java/org/csanchez/jenkins/plugins/kubernetes/volumes\n  volumes:\n  # - type: Secret\n  #   secretName: mysecret\n  #   mountPath: /var/myapp/mysecret\n  NodeSelector: {}\n  # Key Value selectors. Ex:\n  # jenkins-agent: v1\n\nPersistence:\n  Enabled: true\n  ## A manually managed Persistent Volume and Claim\n  ## Requires Persistence.Enabled: true\n  ## If defined, PVC must be created manually before volume will be bound\n  # ExistingClaim:\n\n  ## jenkins data Persistent Volume Storage Class\n  ## If defined, storageClassName: <storageClass>\n  ## If set to \"-\", storageClassName: \"\", which disables dynamic provisioning\n  ## If undefined (the default) or set to null, no storageClassName spec is\n  ##   set, choosing the default provisioner.  (gp2 on AWS, standard on\n  ##   GKE, AWS & OpenStack)\n  ##\n  # StorageClass: \"-\"\n\n  Annotations: {}\n  AccessMode: ReadWriteOnce\n  Size: 8Gi\n  volumes:\n  #  - name: nothing\n  #    emptyDir: {}\n  mounts:\n  #  - mountPath: /var/nothing\n  #    name: nothing\n  #    readOnly: true\n\nNetworkPolicy:\n  # Enable creation of NetworkPolicy resources.\n  Enabled: false\n  # For Kubernetes v1.4, v1.5 and v1.6, use 'extensions/v1beta1'\n  # For Kubernetes v1.7, use 'networking.k8s.io/v1'\n  ApiVersion: extensions/v1beta1\n\n## Install Default RBAC roles and bindings\nrbac:\n  install: false\n  serviceAccountName: default\n  # Role reference\n  roleRef: cluster-admin\n  # Role kind (RoleBinding or ClusterRoleBinding)\n  roleBindingKind: ClusterRoleBinding\n",
		RequiredValues: map[string]string{},
	},
	Spec: SolutionTemplateSpec{
		Requirements: SolutionTemplateSpecRequirements{
			Node: SolutionTemplateSpecRequirementsNode{
				Count: 2,
				CPU: 2,
			},
			App: SolutionTemplateSpecRequirementsApp{
				RBAC: false,
				Namespace: "jenkins",
				ValuesEditRequired: []string{},
				Storage: true,
				LoadBalancer: true,
			},
		},
	},
	Dependencies: SolutionTemplateDependencies{
		Name: "Helm Tiller",
		Value: "helm_tiller",
		Available: false,
		KeysetRequired: false,
		Tag: "latest",
		IsPostBuildCompatible: true,
		IsManagedIndependently: false,
		Dependencies: []string{},
	},
	Version: "0.19.1",
}
