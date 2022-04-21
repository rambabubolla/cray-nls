# Temporary script to deploy argo onto a system for testing
#     NOTE: this will be removed once we have deployment work done
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-creds "admin:admin" docker://quay.io/argoproj/argocli:latest docker://registry.local/argoproj/argocli:latest
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-creds "admin:admin" docker://quay.io/argoproj/argoexec:latest docker://registry.local/argoproj/argoexec:latest
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-creds "admin:admin" docker://quay.io/argoproj/workflow-controller:latest docker://registry.local/argoproj/workflow-controller:latest
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-creds "admin:admin" docker://postgres:12-alpine docker://registry.local/library/postgres:12-alpine
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-creds "admin:admin" docker://minio/minio docker://registry.local/library/minio/minio