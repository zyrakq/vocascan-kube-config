
werf-convert:
  kompose convert -f docker-compose.yml -o ./.helm/templates;
  rm ./.helm/templates/*-persistentvolumeclaim.yaml;
  find ./.helm/templates -type f -exec sed -i "s/'{{{{ \(.*\) }}'/{{{{ \1 }}/g" {} +;
  find ./.helm/templates -type f -exec sed -i "s/\.values/\.Values/g" {} +;

werf-encrypt:
  werf helm secret values encrypt .origin/secret-values.yaml -o .helm/secret-values.yaml
  bash -c 'for filename in .origin/secret/*; do name=${filename##*/}; werf helm secret file encrypt ".origin/secret/$name" -o ".helm/secret/$name"; done;';
werf-decrypt:
  werf helm secret values decrypt .helm/secret-values.yaml -o .origin/secret-values.yaml
  bash -c 'for filename in .helm/secret/*; do name=${filename##*/}; werf helm secret file encrypt ".helm/secret/$name" -o ".origin/secret/$name"; done;';

werf-up-storage:
  kubectl apply -f local-storage.yaml;
  kubectl apply -f vocascandb-pv-0.yaml;
werf-down-storage:
  kubectl delete -f vocascandb-pv-0.yaml;
  kubectl delete -f local-storage.yaml;

werf-up-conf:
  kubectl create namespace vocascan &>/dev/null || exit 0;
  kubectl config set-context --current --namespace=vocascan;
  kubectl apply -Rf './.helm/templates/*-secret.yaml';
werf-down-conf:
  kubectl delete -Rf './.helm/templates/*-secret.yaml';

werf-up *FLAGS:
  werf converge {{FLAGS}};
werf-down *FLAGS:
  werf dismiss {{FLAGS}};
  
werf-clear *FLAGS:
  werf dismiss {{FLAGS}};
  kubectl delete namespace vocascan;
  kubectl delete -f vocascandb-pv-0.yaml;
  kubectl delete -f local-storage.yaml;