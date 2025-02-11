from genericpath import exists
import requests
import os
import json
import git
import time
import logging
import threading

urls = 'https://api.github.com/search/repositories?q=go+bindings&sort=stars&order=desc&per_page=100&page='
repo_path = '/data/github_go/all-repos/'

logger = logging.getLogger()
rq = time.strftime('%Y%m%d%H%M', time.localtime(time.time()))[:-4]
log_file = './get_gobindings_repos.log'
repos = []

def getExistRepo(path):
    exist_repos = []
    dir_or_files = os.listdir(path)
    for f in dir_or_files:
        full_path = os.path.join(path, f)
        if os.path.isdir(full_path):
            exist_repos.append(f)
    return exist_repos

def mkdir(path):
    folder = os.path.exists(path)
    if not folder:                   
        os.makedirs(path)
        logger.info("new folder: %s" % path)

def create_link(src, dst, name):
    src = os.path.join(src, name)
    dst = os.path.join(dst, name)
    try:
        os.symlink(src, dst)
    except Exception as e:
        logger.error(e)
    return

def gitclone(url, to_path):
    try:
        git.Repo.clone_from(url=url, to_path=to_path)
        repos.append(os.path.basename(to_path))
        if not os.path.exists(os.path.join('/data/github_go/gobindings-repos/', os.path.basename(to_path))):
            os.symlink(to_path, os.path.join('/data/github_go/gobindings-repos/', os.path.basename(to_path)))
    except Exception as e:
        logger.error(e)
    return 

def main():
    all_repos = getExistRepo('/data/github_go/all-repos/')

    i = 1
    logger.info('now getting page 1')
    # 获取搜索结果
    try:
        r = requests.get(urls+str(i))
        result = json.loads(r.content.decode('UTF-8'))
    except Exception as e:
        logger.error(e)
        return
    
    # count = 0
    # 尝试clone每一个库
    for repo in result['items']:
        name = repo['name']
        if name not in all_repos:
            mkdir('/data/github_go/all-repos-info/'+name)
            with open('/data/github_go/all-repos-info/'+name+'/'+name+'.json','w+') as f:
                json.dump(repo,f,indent=4)
                f.close()
            create_link('/data/github_go/all-repos-info/', '/data/github_go/gobindings-repos-info/', name)
            
            clone_url = repo['clone_url']
            
            # git clone
            logger.info('now getting repo:' + repo['full_name'] + ' from ' + clone_url + ' to ' + repo_path+name)
            # if count % 4 == 0:
            #     t0 = threading.Thread(target=gitclone, args=(clone_url, repo_path+name))
            #     t0.start()
            # elif count % 4 == 1:
            #     t1 = threading.Thread(target=gitclone, args=(clone_url, repo_path+name))
            #     t1.start()
            # elif count % 4 == 2:
            #     t2 = threading.Thread(target=gitclone, args=(clone_url, repo_path+name))
            #     t2.start()
            # else:
            gitclone(clone_url, repo_path+name)          
            # count += 1

    repo_list = {}
    logging.info('Clone over! Total repos num:'+str(len(repos)))
    json_file = open('gobindings_repo.json','w+')
    repo_list['repo_num'] = len(repos)
    repo_list['repos'] = repos
    json.dump(repo_list,json_file,indent=4)
    json_file.close()

if __name__ == '__main__':
    logger.setLevel(logging.INFO)
    fh = logging.FileHandler(log_file, mode='a')
    formatter = logging.Formatter("%(asctime)s - %(filename)s[line:%(lineno)d] - %(levelname)s: %(message)s")
    fh.setFormatter(formatter)
    logger.addHandler(fh)
    main()
