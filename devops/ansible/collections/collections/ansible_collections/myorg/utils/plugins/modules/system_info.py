#!/usr/bin/python

from __future__ import absolute_import, division, print_function
__metaclass__ = type

DOCUMENTATION = r"""
---
module: system_info
short_description: Returns basic system information
description:
  - A custom module bundled in the myorg.utils collection.
  - Demonstrates how to write and distribute custom modules via collections.
options:
  message:
    description: A custom message to include in the output.
    type: str
    default: "Hello from myorg.utils collection!"
"""

EXAMPLES = r"""
- name: Gather system info
  myorg.utils.system_info:
    message: "Deployed via Ansible Collections!"
  register: info

- ansible.builtin.debug:
    msg: "{{ info.hostname }} running Python {{ info.python_version }}"
"""

RETURN = r"""
message:
  description: The message passed as input.
  type: str
hostname:
  description: The system hostname.
  type: str
platform:
  description: The OS platform name.
  type: str
python_version:
  description: The Python interpreter version.
  type: str
"""

import platform
import socket

from ansible.module_utils.basic import AnsibleModule


def main():
    module = AnsibleModule(
        argument_spec=dict(
            message=dict(type="str", default="Hello from myorg.utils collection!"),
        )
    )

    result = dict(
        changed=False,
        message=module.params["message"],
        hostname=socket.gethostname(),
        platform=platform.system(),
        python_version=platform.python_version(),
    )

    module.exit_json(**result)


if __name__ == "__main__":
    main()
