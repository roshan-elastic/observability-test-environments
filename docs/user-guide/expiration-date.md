# Cluster expiration date

The cluster expiration date is the date when the cluster will be automatically deleted. The cluster expiration date is set when the cluster is created. The cluster expiration date can be extended by the cluster owner.
Each template has a different expiration time, the most common expiration time is 30 days.
For serverless clusters, the expiration date is set to 3 day.
For the clusters created by automation (CI/CD), the expiration date is set to 1 day.

You can change the expiration date of your cluster by using the following command:

```bash
UPDATE_DATE_ISO='{"expire_date":"2024-01-31"}'
CLUSTER_NAME=my-cluster

oblt-cli cluster update --cluster-name ${CLUSTER_NAME} --parameters ${UPDATE_DATE_ISO}
```

!!! Warning

    SNAPSHOT deployments use a trial license that expires after 30 days. If you want to extend the expiration date of your cluster, [you have to upload a new license](./push-licenses.md).
