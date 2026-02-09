#!/usr/bin/env python3
import os
import logging
from pathlib import Path
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler
from kubernetes import client, config

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


def parse_properties_file(file_path: str) -> dict:
    props = {}
    try:
        with open(file_path, 'r') as f:
            for line in f:
                line = line.strip()
                if line and not line.startswith('#') and '=' in line:
                    key, value = line.split('=', 1)
                    props[key.strip()] = value.strip()
        return props
    except Exception as e:
        logger.error(f"Failed to parse {file_path}: {e}")
        return {}


def sync_properties(watch_dir: str, namespace: str, configmap_name: str):
    all_props = {}
    
    for file_path in Path(watch_dir).rglob('*.properties'):
        props = parse_properties_file(str(file_path))
        all_props.update(props)
        logger.debug(f"Loaded {len(props)} properties from {file_path}")
    
    if not all_props:
        logger.debug("No properties to sync")
        return
    
    v1 = client.CoreV1Api()
    try:
        cm = v1.read_namespaced_config_map(configmap_name, namespace)
        cm.data = all_props
        v1.replace_namespaced_config_map(configmap_name, namespace, cm)
        logger.info(f"Updated ConfigMap '{configmap_name}' with {len(all_props)} entries")
    except client.exceptions.ApiException as e:
        if e.status == 404:
            cm = client.V1ConfigMap(
                metadata=client.V1ObjectMeta(name=configmap_name, namespace=namespace),
                data=all_props
            )
            v1.create_namespaced_config_map(namespace, cm)
            logger.info(f"Created ConfigMap '{configmap_name}' with {len(all_props)} entries")
        else:
            logger.error(f"Failed to update ConfigMap: {e}")


class PropertiesFileHandler(FileSystemEventHandler):
    def __init__(self, watch_dir: str, namespace: str, configmap_name: str):
        self.watch_dir = watch_dir
        self.namespace = namespace
        self.configmap_name = configmap_name
    
    def on_modified(self, event):
        if not event.is_directory and event.src_path.endswith('.properties'):
            logger.info(f"File change detected: {event.src_path}")
            sync_properties(self.watch_dir, self.namespace, self.configmap_name)
    
    def on_created(self, event):
        if not event.is_directory and event.src_path.endswith('.properties'):
            logger.info(f"File created: {event.src_path}")
            sync_properties(self.watch_dir, self.namespace, self.configmap_name)


def main():
    watch_dir = os.getenv('WATCH_DIR', '/etc/config')
    namespace = os.getenv('POD_NAMESPACE', 'default')
    configmap_name = os.getenv('CONFIGMAP_NAME', 'properties')
    
    logger.info(f"Starting Python sidecar")
    logger.info(f"Watch dir: {watch_dir}")
    logger.info(f"Namespace: {namespace}")
    logger.info(f"ConfigMap: {configmap_name}")
    
    load_kube_config()
    
    sync_properties(watch_dir, namespace, configmap_name)
    
    event_handler = PropertiesFileHandler(watch_dir, namespace, configmap_name)
    observer = Observer()
    observer.schedule(event_handler, watch_dir, recursive=True)
    observer.start()
    
    logger.info("Watcher started, monitoring for changes...")
    
    try:
        observer.join()
    except KeyboardInterrupt:
        observer.stop()
        logger.info("Shutting down...")
        observer.join()


if __name__ == "__main__":
    main()
