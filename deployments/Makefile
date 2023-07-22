.PHONY: render install upgrade delete manifest logs

render:
	helm template snowflake-was . -f values.yaml --namespace demo --debug > rendered.yaml

install:
	helm install snowflake-was . --namespace snowflake-was --create-namespace --values ./values.yaml

upgrade:
	helm upgrade snowflake-was . --namespace snowflake-was --values ./values.yaml

delete:
	helm uninstall snowflake-was --namespace snowflake-was

manifest:
	helm get manifest snowflake-was --namespace snowflake-was > manifest.yaml

logs:
	kubectl logs --namespace snowflake-was --selector release=snowflake-was
