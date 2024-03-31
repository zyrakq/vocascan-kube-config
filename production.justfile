_synth-production:
  cd ./.cdk8s && go run ./src -e=Production

_encrypt-production:
  werf helm secret values encrypt .cdk8s/secret-values.production.yaml -o .helm-production/secret-values.yaml
  bash -c 'werf helm secret file encrypt ".cdk8s/vocascan.config.production.js" -o ".helm-production/secret/vocascan.config.js"';
_decrypt-production:
  werf helm secret values decrypt .helm-production/secret-values.yaml -o .cdk8s/secret-values.production.yaml
  bash -c 'werf helm secret file decrypt ".helm-production/secret/vocascan.config.js" -o ".cdk8s/vocascan.config.production.js"';


_up-production *FLAGS:
  werf converge --config='werf.production.yaml' {{FLAGS}};
_down-production *FLAGS:
  werf dismiss --config='werf.production.yaml' {{FLAGS}};
