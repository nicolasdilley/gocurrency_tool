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
   p+=1
   payload = {'q':'language:go', 'sort':'stars', 'order':'desc','per_page':'100','page':str(p)}
   r = s.get('https://api.github.com/search/repositories', params=payload)
   data = r.json()
   # print(data)
   for repo in data['items']:
       i+=1
       r = s.get('https://api.github.com/repos/'+str(repo['full_name']))
       repodata = r.json()
       if len(repodata['topics']) == 0:
          print(str(repo['full_name']),','
                # , str(repo['description']),','
                , str(repo['watchers_count']))
       else:
          print(str(repo['full_name']),','
                # , str(repo['description']),','
                , str(repo['watchers_count']), ','
                , ','.join(repodata['topics']))
          # keywords = keywords+(repodata['topics'])
          
       if int(str(repo['watchers_count'])) < 1000:
           cont = False;
    
print(str(i)," projects found")
# print(','.join(set(keywords)))
