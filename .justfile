
werf-convert:
  kompose convert -f docker-compose.yml -o ./.helm/templates;
  rm ./.helm/templates/*-persistentvolumeclaim.yaml;
  find ./.helm/templates -type f -exec sed -i "s/'{{{{ \(.*\) }}'/{{{{ \1 }}/g" {} +;
  mv .helm/templates/*-secret.yaml .kube/secret;

werf-encrypt:
  werf helm secret values encrypt .raw/secret-values.yaml -o .helm/secret-values.yaml
  bash -c 'for filename in .raw/secret/*; do name=${filename##*/}; werf helm secret file encrypt ".raw/secret/$name" -o ".helm/secret/$name"; done;';
werf-decrypt:
  werf helm secret values decrypt .helm/secret-values.yaml -o .raw/secret-values.yaml
  bash -c 'for filename in .helm/secret/*; do name=${filename##*/}; werf helm secret file encrypt ".helm/secret/$name" -o ".raw/secret/$name"; done;';

werf-up-storage:
  kubectl apply -f local-storage.yaml;
  kubectl apply -f vocascandb-pv-0.yaml;
werf-down-storage:
  kubectl delete -f vocascandb-pv-0.yaml;
  kubectl delete -f local-storage.yaml;

werf-up-conf:
  kubectl create namespace vocascan &>/dev/null || exit 0;
  kubectl config set-context --current --namespace=vocascan;
  kubectl apply -Rf ./.kube/secret/;
werf-down-conf:
  kubectl apply -Rf ./.kube/secret/;

werf-up:
  werf converge;
werf-down:
  werf dismiss;
  
werf-clear:
  werf dismiss;
  kubectl delete namespace vocascan;
  kubectl delete -f vocascandb-pv-0.yaml;
  kubectl delete -f local-storage.yaml;