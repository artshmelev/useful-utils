#!/usr/bin/env python2.7

from jira import JIRA
import getpass
import warnings
warnings.filterwarnings('ignore')

if __name__ == '__main__':
    jira = JIRA(
        options={'server': 'https://jira'},
        basic_auth=('user', getpass.getpass('password> '))
    )
    lst = ['('+i.key+') '+i.fields.summary for i in jira.search_issues('query', maxResults=6)]
    for i in lst:
        print i
