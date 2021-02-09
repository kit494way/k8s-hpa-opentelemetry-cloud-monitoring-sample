# Sample of Kubernetes HorizontalPodAutoScaler, OpenTelemetry and Cloud Monitoring

## クラスタ作成

```sh
$ gcloud container clusters create otel-example
```

## custom-metric-stackdriver-adapter をデプロイ

必要なロールを作成する権限を自身のアカウントに付与する。

```sh
$ kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin --user "$(gcloud config get-value account)"
```

custom-metric-stackdriver-adapter をデプロイする。

```sh
$ kubectl apply -f https://raw.githubusercontent.com/GoogleCloudPlatform/k8s-stackdriver/master/custom-metrics-stackdriver-adapter/deploy/production/adapter_new_resource_model.yaml
```

## OpenTelemetry Collector と stackdriverexporter をデプロイ

manifests/otel-collector.yaml の `$PROJECT_ID` を GCP のプロジェクト ID で置き換える。

```sh
$ kubectl apply -f manifests/otel-collector.yaml
```

## メトリクスを送信する Docker Image をビルド

```sh
$ cd app
$ gcloud builds submit --tag gcr.io/$PROJECT_ID/custom-metric
```

Google Cloud Registry にイメージがプッシュされる。

## カスタムメトリクスを送信するアプリをデプロイ

manifests/custom-metric.yaml の `$PROJECT_ID` を GCP のプロジェクト ID で置換する。

```sh
$ kubectl apply -f manifests/custom-metric.yaml
```

## HorizontalPodAutoScaler をデプロイ

```sh
$ kubectl apply -f manifests/custom-metric-hpa.yaml
```
