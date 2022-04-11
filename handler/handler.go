package handler

import (
  "context"
  "encoding/json"
  "fmt"
  "net/http"
 // "reflect"
 // "strings"

  "github.com/go-logr/logr"
  //metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

  corev1 "k8s.io/api/core/v1"
  "k8s.io/apimachinery/pkg/labels"
  "k8s.io/apimachinery/pkg/selection"
  //utilerrors "k8s.io/apimachinery/pkg/util/errors"
  "sigs.k8s.io/controller-runtime/pkg/client"
  "sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
  gitOpsAppLabel = "policygenerator.admission.kubernetes.io"
)

// +kubebuilder:webhook:path=/mutate,mutating=true,failurePolicy=ignore,groups="",resources=pods,verbs=create,versions=v1,name=mpod.redhatcop.redhat.io,sideEffects=None,admissionReviewVersions={v1,v1beta1}
// +kubebuilder:rbac:groups=redhatcop.redhat.io,resources=podpresets,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;update;patch

// PolicyGeneratorMutator mutates Pods
type PolicyGeneratorMutator struct {
  Client  client.Client
  decoder *admission.Decoder
  Log     logr.Logger
}

// PolicyGeneratorMutator adds an annotation to every incoming pods.
func (a *PolicyGeneratorMutator) Handle(ctx context.Context, req admission.Request) admission.Response {
  logger := a.Log.WithValues("policy-generator-webhook", fmt.Sprintf("%s/%s", req.Namespace, req.Name))

  // Ignore all calls to subresources or resources other than pods.
  // Ignore all operations other than CREATE.
  if len(req.SubResource) != 0 || req.Resource.Group != "" || req.Operation != "CREATE" {
    return admission.Allowed("")
  }

  logger.Info("test logger")
  pod := &corev1.Pod{}

  err := a.decoder.Decode(req, pod)
  if err != nil {
    return admission.Errored(http.StatusBadRequest, err)
  }

  // Begin Mutation

  if _, isMirrorPod := pod.Annotations[corev1.MirrorPodAnnotationKey]; isMirrorPod {
    return admission.Allowed("Mirror Pod")
  }

  // Ignore if not the gitops repo server
  gitOpsReq, _ := labels.NewRequirement("app.kubernetes.io/name", selection.Equals, []string{"openshift-gitops-repo-server"})
 
  // Init and add to selector.
  selector := labels.NewSelector()
  selector = selector.Add(*gitOpsReq)
  // check if the pod labels match the selector
  if !selector.Matches(labels.Set(pod.Labels)) {
    return admission.Allowed("Exclude pod not gitops repo server")
  }

  // Create the emptyDir to hold the generator
  emptyVolume := corev1.Volume{
                    Name: "policy-generator",
                    VolumeSource: corev1.VolumeSource{
                      EmptyDir: &corev1.EmptyDirVolumeSource{},
                    },
                  }

  // Create the volume mount
  emptyMount := corev1.VolumeMount{
                    Name:      "policy-generator",
                    MountPath: "/etc/kustomize",
                }

  // Create the initContainer to copy the generator into the emptyDir

  // Add the emptyDir volume to the pod
  pod.Spec.Volumes = append(pod.Spec.Volumes, emptyVolume)

  // Add the intiContainer to the pod

  // Add the volumemount to the container(s)
  vms := pod.Spec.Containers[0].VolumeMounts
  pod.Spec.Containers[0].VolumeMounts = append(vms, emptyMount)
  

  // End Mutation
  marshaledPod, err := json.Marshal(pod)
  if err != nil {
    return admission.Errored(http.StatusInternalServerError, err)
  }

  return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

// PodPresetMutator implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (a *PolicyGeneratorMutator) InjectDecoder(d *admission.Decoder) error {
  a.decoder = d
  return nil
}




// // applyPodPresetsOnPod updates the PodSpec with merged information from all the
// // applicable PodPresets. It ignores the errors of merge functions because merge
// // errors have already been checked in safeToApplyPodPresetsOnPod function.
// func applyPodPresetsOnPod(pod *corev1.Pod, podPresets []*redhatcopv1alpha1.PodPreset) {
//   if len(podPresets) == 0 {
//     return
//   }

//   volumes, _ := mergeVolumes(pod.Spec.Volumes, podPresets)
//   pod.Spec.Volumes = volumes

//   for i, ctr := range pod.Spec.Containers {
//     applyPodPresetsOnContainer(&ctr, podPresets)
//     pod.Spec.Containers[i] = ctr
//   }
//   for i, iCtr := range pod.Spec.InitContainers {
//     applyPodPresetsOnContainer(&iCtr, podPresets)
//     pod.Spec.InitContainers[i] = iCtr
//   }

//   // add annotation
//   if pod.ObjectMeta.Annotations == nil {
//     pod.ObjectMeta.Annotations = map[string]string{}
//   }

//   for _, pp := range podPresets {
//     pod.ObjectMeta.Annotations[fmt.Sprintf("%s/podpreset-%s", annotationPrefix, pp.GetName())] = pp.GetResourceVersion()
//   }
// }

// // applyPodPresetsOnContainer injects envVars, VolumeMounts and envFrom from
// // given podPresets in to the given container. It ignores conflict errors
// // because it assumes those have been checked already by the caller.
// func applyPodPresetsOnContainer(ctr *corev1.Container, podPresets []*redhatcopv1alpha1.PodPreset) {
//   envVars, _ := mergeEnv(ctr.Env, podPresets)
//   ctr.Env = envVars

//   volumeMounts, _ := mergeVolumeMounts(ctr.VolumeMounts, podPresets)
//   ctr.VolumeMounts = volumeMounts

//   envFrom, _ := mergeEnvFrom(ctr.EnvFrom, podPresets)
//   ctr.EnvFrom = envFrom
// }
