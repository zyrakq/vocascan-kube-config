_synth-development:
  cd ./.cdk8s && go run ./src -e=Development

_encrypt-development:
  werf helm secret values encrypt .cdk8s/secret-values.development.yaml -o .helm-development/secret-values.yaml
  bash -c 'werf helm secret file encrypt ".cdk8s/vocascan.config.development.js" -o ".helm-development/secret/vocascan.config.js"';
_decrypt-development:
  werf helm secret values decrypt .helm-development/secret-values.yaml -o .cdk8s/secret-values.development.yaml
  bash -c 'werf helm secret file decrypt ".helm-development/secret/vocascan.config.js" -o ".cdk8s/vocascan.config.development.js"';


_up-development *FLAGS:
  werf converge --config='werf.development.yaml' {{FLAGS}};
_down-development *FLAGS:
  werf dismiss --config='werf.development.yaml' {{FLAGS}};
