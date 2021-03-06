# **Change Log**

This project adheres to [Semantic Versioning](http://semver.org/). Additionally, below are the change categories that may be associated with each release version.

- **Added** for new features.
- **Changed** for changes in existing functionality.
- **Deprecated** for once-stable features removed in upcoming releases.
- **Removed** for deprecated features removed in this release.
- **Fixed** for any bug fixes.
- **Security** for any security changes or fixes for vulnerabilities.

### **[0.0.7]**

#### Fixed

*  Enable unit tests for all environment-operator packages. [[BITE-1472](https://agile-jira.pearson.com/browse/BITE-1472)]

### **[0.0.6] - 2017-09-13 [RELEASED]**

#### Fixed

*  Ensure k8s resources are only applied if a deployment is made for that Bitesize Service. [[BITE-1634] (https://agile-jira.pearson.com/browse/BITE-1634)]

### **[0.0.5] - 2017-09-06 [RELEASED]**

#### Fixed

* Bug caused by annotations with pods continuously upgrading.

#### Changed

* Service creation logic has changed. Now kubernetes resource will only be created after the deployment fact (i.e. we will not create service, ingress etc. resources for the service that is not yet deployed as a pod)
* (Internals) Pod logs are no longer a part of bitesize environment object.

### **[0.0.4] - 2017-09-01 [RELEASED]**

#### Added

*  Support for kubernetes service annotations. [[BITE-1511](https://agile-jira.pearson.com/browse/BITE-1511)]

### **[0.0.3] - 2017-08-31 [RELEASED]**

#### Added

*  Support for configuring horizontal pod autoscaling. [[BITE-1433](https://agile-jira.pearson.com/browse/BITE-1433)]
*  Added new environment operator endpoint for Pod Status. [[BITE-1484](https://agile-jira.pearson.com/browse/BITE-1484)]
*  Custom Docker registry support added for pod spec. [[BITE-1448](https://agile-jira.pearson.com/browse/BITE-1448)]
*  Environment Operator build/release pipeline now managed by TravisCI. [[BITE-1473](https://agile-jira.pearson.com/browse/BITE-1473)]
*  Add error handling for secrets defined in environment.bitesize files for deployments. [[BITE-1465](https://agile-jira.pearson.com/browse/BITE-1465)]

### **[0.0.2] - 2017-01-17 [RELEASED]**

* Validator command added for validation of environment.bitesize file

### **[0.0.1] - 2017-01-17 [RELEASED]**

* Original release of environment operator.
