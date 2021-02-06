# dynamic-update-configmap-env-vars
Demo on how to keep pod's environment variables in sync with any updates in the corresponding ConfigMap

### Command to test the service
kubectl run curl-test --image=radial/busyboxplus:curl -i --tty --rm