# -*- mode: Python -*-

def helmfile():
  return local("helmfile -e local -f helmfile.d/01-app.yaml template")

load('ext://nerdctl', 'nerdctl_build')
load('ext://namespace', 'namespace_create', 'namespace_inject')

ns = 'exporter-weather-local-v1'
nerdctl_build(
    ref='exporter-weather',
    context='.',
)
helm_yaml = helmfile()
k8s_yaml(helm_yaml)

namespace_create(ns)
helm_yaml = helmfile()
k8s_yaml(namespace_inject(helm_yaml, ns))