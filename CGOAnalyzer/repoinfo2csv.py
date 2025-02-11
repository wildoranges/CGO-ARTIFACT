import json
import os
import datetime
import psycopg2
import shutil

infos_path = '/data/github_go/all-repos-info/'
repo_path = '/data/github_go/all-repos/'

key2json = {
    'repo_name': 'name',
    'url': 'html_url',
    'stars': 'stargazers_count',
    'loc':None,
    'size': 'size',
    'forks_count': 'forks_count',
    'issues': 'open_issues_count',
    'created_at': 'created_at',
    'updated_at': 'updated_at',
    'pushed_at': 'pushed_at',
    'description':'description',
    'archived':'archived'
}

date_item = {
    'updated_at',
    'created_at'
}

CPP_SUFFIX_SET = {'.h', '.hpp', '.hxx', '.c', '.cpp', '.cc', '.cxx', '.hh'}
GO_SUFFIX_SET = {'.go'}

EXCLUDE_DIR_SET = {'vendor','test','bin'}

cpp_lines = 0
go_lines = 0
total_lines = 0


def get_loc(repo_name):
    global CPP_SUFFIX_SET, GO_SUFFIX_SET
    global cpp_lines,go_lines,total_lines,repo_path
    cpp_lines = 0
    go_lines = 0
    total_lines = 0
    cur_path = os.path.join(repo_path, repo_name)
    is_exists = os.path.exists(cur_path)
    if not is_exists:
        return None
    list_files(cur_path)
    total_lines = cpp_lines + go_lines
    return total_lines


def count_lines(path):
    global CPP_SUFFIX_SET, GO_SUFFIX_SET
    global cpp_lines,go_lines,total_lines
    suffix = os.path.splitext(path)[-1]
    cnt = 0
    if (suffix in CPP_SUFFIX_SET) or (suffix in GO_SUFFIX_SET):
        # print("now checking file: {}".format(path))
        f = open(path,'rb')
        last_data = '\n'
        while True:
            data = f.read(0x400000)
            if not data:
                break
            cnt += data.count(b'\n')        
            last_data = data
        if last_data[-1:] != b'\n':
            cnt += 1
        f.close()
    
    if suffix in CPP_SUFFIX_SET:
        cpp_lines += cnt
    elif suffix in GO_SUFFIX_SET:
        go_lines += cnt


def list_files(path):
    files = os.listdir(path)
    for f in files:
        cur_path = os.path.join(path,f)
        if os.path.isfile(cur_path):
            if not os.path.islink(cur_path):
                count_lines(cur_path)
        if os.path.isdir(cur_path):
            if f not in EXCLUDE_DIR_SET and not os.path.islink(cur_path):
                list_files(cur_path)

def insert2table(cur, conn, data):
    ptn = "INSERT INTO repository (repo_name, url, stars, loc, size, forks_count, issues, created_at, updated_at, repo_type) \
        VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s)"
    cur.execute(ptn, (data['repo_name'], data['url'], data['stars'], data['loc'], data['size'], data['forks_count'], data['issues'], data['created_at'], data['updated_at'], data['repo_type']))
    conn.commit()

def get_exist_repos(path):
    exist_repos = []
    dir_or_files = os.listdir(path)
    for f in dir_or_files:
        full_path = os.path.join(path, f)
        if os.path.isdir(full_path):
            exist_repos.append(f)
    return exist_repos

def main():
    # conn = psycopg2.connect(database="cgo", user="postgres", password="s4plususer", host="127.0.0.1", port="5432")
    # cur = conn.cursor()

    exist_repos = get_exist_repos('/data/github_go/all-repos/')
    binding_repos = get_exist_repos('/data/github_go/gobindings-repos-info/')
    dataset = []
    for repo in exist_repos:
        info_path = os.path.join(infos_path, repo)
        info_path = os.path.join(info_path, repo + '.json')
        f = open(info_path, 'r')
        repo_info = json.load(f)
        data = {}

        for key in key2json.keys():
            if key2json[key]:
                if key2json[key] in date_item:
                    day = datetime.datetime.strptime(repo_info[key2json[key]],"%Y-%m-%dT%H:%M:%SZ")
                    date = datetime.datetime.strftime(day,"%Y-%m-%d %H:%M:%S")
                    data[key] = date
                else:
                    data[key] = repo_info[key2json[key]]
            else:
                data[key] = get_loc(repo)
        data['repo_type'] = ""
        data['bindings'] = (repo in binding_repos)
        dataset.append(data)
        f.close()

        # insert2table(cur, conn, data)
    print('start to dump result')
    f = open('all-repos-info.json', 'w+')
    json.dump(dataset, f)
    f.close()
    # conn.close()


if __name__ == '__main__':
    main()
    # exist_repos = get_exist_repos('/data/github_go/all-repos/')
    # for repo in exist_repos:
    #     # repo_path = os.path.join('/data/github_go/all-repos/', repo)
    #     loc = get_loc(repo)
    #     if loc == 0:
    #         # print(repo)
    #         print(os.path.join(repo_path, repo))
    #         shutil.rmtree(os.path.join(repo_path, repo))
    #     # print(loc)