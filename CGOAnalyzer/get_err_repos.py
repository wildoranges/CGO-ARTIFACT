import os
import git
import json
import threading

repo_info_base = '/data/github_go/all-repos-info/'

def gitclone(url, to_path):
    try:
        print("try to clone " + url + " to " + to_path)
        git.Repo.clone_from(url=url, to_path=to_path)
    except Exception as e:
        print("error: " + to_path)
        print(e)

def get_error_repos():
    f1 = open('err_repos.txt', 'r')
    lines = f1.readlines()
    # count = 0
    for line in lines:
        line = line.strip('\n')
        repo = os.path.basename(line)
        repo_info_path = os.path.join(repo_info_base, repo)
        repo_info_path = os.path.join(repo_info_path, repo+'.json')
        f = open(repo_info_path, 'r')
        repo_info = json.load(f)
        url = repo_info['clone_url']

        gitclone(url, line)

        f.close()
        # count += 1

    f1.close()

if __name__ == '__main__':
    get_error_repos()