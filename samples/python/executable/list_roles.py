import boto3
from tabulate import tabulate
import sys

def list_iam_roles():
    try:
        iam = boto3.client('iam')
        paginator = iam.get_paginator('list_roles')
        
        roles_data = []
        for page in paginator.paginate():
            for role in page['Roles']:
                roles_data.append([
                    role['RoleName'],
                    role['CreateDate'].strftime('%Y-%m-%d'),
                    role['Path']
                ])
        
        if roles_data:
            print(tabulate(roles_data, headers=['Role Name', 'Created Date', 'Path'], tablefmt='grid'))
            print(f"\nTotal roles: {len(roles_data)}")
        else:
            print("No IAM roles found.")
            
    except Exception as e:
        print(f"Error: {str(e)}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    list_iam_roles()
