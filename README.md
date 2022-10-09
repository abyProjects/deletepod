# deletepod
This is a commad line utility to delete nginx standalone pod. <br>
User have to give the name of pod, namespce of the pod and
the token for authentication.

By default this CLI supports the deletion of "nginx" pod within the namespace "test".<br>
Default token for authentication is "abcd"<br><br>
NOTE: This is basic authentication check not a secure one. You can either pass the token with value or with the filepath which contains the value.

Please follow the instructions below to run the deletepod application

**non-cluster cli execution step**
1. go run main.go pod --name nginx --namespace test --token "path-of-tokenfile"
2. go run main.go pod --name nginx --namespace test --token abcd

**to deploy cli in a kubernetes cluster**
1. Edit the Dockerfile and provide the ENV variables with respective PODNAME, NAMESPACE and TOKEN <br>
NOTE: by default this application is designed to delete "nginx" pod within the namespace "test" using token "abcd"
2. After edit is completed, build the Dockerfile with tag "aby0516/deletepod:v2"[NOTE: "optional"]<br>
    - "docker build -t aby0516/deletepod:v2 ."
3. Push the docker image to docker hub. [NOTE: "optional" unless user made a change, by default image is already build and pushed to docker ]
    - "docker push aby0516/deletepod:v2"
4. To create namespace and to deploy nginx pod, clusterRole and deletepod application. 
   Move to yamlFiles directory under the project path and execute the apply command as shown below.
    - "cd yamlFiles"
    - "kubectl apply -f testNamespace.yaml -f clusterRole.yaml  -f nginxpod.yaml -f deletepoDepoly.yaml"

The deletepod application will automatically find and delete nginx pod if deployed under namespace test.