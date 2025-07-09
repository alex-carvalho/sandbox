import os

def main():
    dir_name = "example_dir"
    file_name = f"{dir_name}/example_file.txt"
    hard_link = f"{dir_name}/example_file_hardlink.txt"
    sym_link = f"{dir_name}/example_file_symlink.txt"

    # Create directory
    try:
        os.mkdir(dir_name, 0o755)
        print(f"Directory '{dir_name}' created.")
    except FileExistsError:
        print(f"Directory '{dir_name}' already exists.")
    except Exception as e:
        print(f"mkdir error: {e}")

    # Create a file inside the directory
    try:
        with open(file_name, "w") as f:
            f.write("Hello, file system syscalls!\n")
        print(f"File '{file_name}' created.")
    except Exception as e:
        print(f"open error: {e}")

    # Create a hard link
    try:
        os.link(file_name, hard_link)
        print(f"Hard link '{hard_link}' created.")
    except Exception as e:
        print(f"link error: {e}")

    # Create a symbolic link
    try:
        os.symlink("example_file.txt", sym_link)
        print(f"Symbolic link '{sym_link}' created (points to 'example_file.txt').")
    except Exception as e:
        print(f"symlink error: {e}")

    # Change file permissions
    try:
        os.chmod(file_name, 0o600)
        print(f"Permissions of '{file_name}' changed to 0600.")
    except Exception as e:
        print(f"chmod error: {e}")

    # Change file ownership (to current user)
    try:
        os.chown(file_name, os.getuid(), os.getgid())
        print(f"Ownership of '{file_name}' changed to current user.")
    except Exception as e:
        print(f"chown error: {e}")

    # Remove all files and links inside the directory before removing the directory itself
    for path, desc in [(file_name, "file"), (hard_link, "hard link"), (sym_link, "symbolic link")]:
        try:
            os.unlink(path)
            print(f"{desc.capitalize()} '{path}' removed.")
        except Exception as e:
            print(f"unlink {desc} error: {e}")

    # Remove the directory itself
    try:
        os.rmdir(dir_name)
        print(f"Directory '{dir_name}' removed.")
    except Exception as e:
        print(f"rmdir error: {e}")

if __name__ == "__main__":
    main()