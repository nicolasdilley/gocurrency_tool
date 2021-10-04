import json
import requests
import getpass
from pprint import pprint

# curl "https://api.github.com/search/repositories?q=language:go&sort=stars&order=desc" > gorepo.json


i=0
p=0
cont = True

username = input("GitHub username: ")
pswd = getpass.getpass()

s = requests.Session()
s.auth = (username, pswd)
s.headers.update({'Accept':'application/vnd.github.mercy-preview+json'})


keywords = []

while cont:
    payload = {'q':'language:go', 'sort':'stars', 'order':'desc','per_page':'150','page':str(p)}
    r = s.get('https://api.github.com/search/repositories', params=payload)
    data = r.json()
    # print(data)
    for repo in data['items']:
       i+=1
       r = s.get('https://api.github.com/repos/'+str(repo['full_name']))
       repodata = r.json()
      
       print(str(repo['full_name']),','
                # , str(repo['description']),','
                , str(repo['watchers_count']))

         # keywords = keywords+(repodata['topics'])
          
    if i > 150:
        cont = False;
    
print(str(i)," projects found")
# print(','.join(set(keywords)))
