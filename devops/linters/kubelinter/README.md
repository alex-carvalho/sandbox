## kubelinter

https://docs.kubelinter.io/#/


```shell
# run linter on a file
kube-linter lint pod.yaml

# check the list of rules
kube-linter checks list

# is possible to add anotations on file to ignore some rules, I don't like this idea
#  annotations:
#    ignore-check.kube-linter.io/unset-cpu-requirements : "cpu requirements not required"


# pass config file to linter
kube-linter lint replica-set.yaml --config kubelinter-config.yaml
```