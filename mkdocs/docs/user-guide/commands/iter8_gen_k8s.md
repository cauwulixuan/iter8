---
template: main.html
title: "Iter8 Gen K8s"
hide:
- toc
---

## iter8 gen k8s

Generate manifest for running experiment in Kubernetes

```
iter8 gen k8s [flags]
```

### Examples

```

# Generate Kubernetes manifest
iter8 gen k8s
```

### Options

```
  -a, --app string   app to be associated with an experiment, default is 'default' (default "default")
  -h, --help         help for k8s
  -i, --id string    if not specified, a randomly generated identifier will be used
```

### Options inherited from parent commands

```
      --set stringArray          set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)
      --set-file stringArray     set values from respective files specified via the command line (can specify multiple or separate values with commas: key1=path1,key2=path2)
      --set-string stringArray   set STRING values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)
  -f, --values strings           specify values in a YAML file or a URL (can specify multiple)
```

### SEE ALSO

* [iter8 gen](iter8_gen.md)	 - Render templates with values

###### Auto generated by spf13/cobra on 13-Jan-2022