#!/usr/bin/env python3
import logging
import os
from kubernetes import client, config, watch

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


def load_kube_config():
    try:
        config.load_incluster_config()
        logger.info("Using in-cluster configuration")
    except config.config_exception.ConfigException:
        config.load_kube_config()
        logger.info("Using local kubeconfig")


def print_configmap_data(cm_data: dict):
    if cm_data:
        logger.info(f"ConfigMap data ({len(cm_data)} entries):")
        for key, value in cm_data.items():
            logger.info(f"  {key} = {value}")
    else:
        logger.warning("ConfigMap is empty")


def watch_configmap(namespace: str, name: str):
    v1 = client.CoreV1Api()
    w = watch.Watch()

    logger.info(f"Starting ConfigMap watcher: namespace={namespace}, configmap={name}")

    try:
        for event in w.stream(v1.list_namespaced_config_map, namespace, field_selector=f"metadata.name={name}"):
            cm = event['object']
            event_type = event['type']
            logger.info(f"ConfigMap {event_type}")
            print_configmap_data(cm.data)
    except Exception as e:
        logger.error(f"Watch error: {e}")
        w.stop()


def main():
    load_kube_config()

    namespace = os.getenv("POD_NAMESPACE", "default")
    configmap_name = os.getenv("CONFIGMAP_NAME", "properties")

    watch_configmap(namespace, configmap_name)


if __name__ == "__main__":
    main()


if __name__ == "__main__":
    main()
