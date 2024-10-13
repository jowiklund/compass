# Compass
MVP:
[x] Log in with account

[x] Get zones

[x] Chose what zones to create permissions for

[x] Show a permission creator for each zone where the user can select what permissions to use for the zone
    - Use a stack

[x] When the user has created their permissions, show a form that lets the user complete their token request with name, validity and password
    - Use the same name for the token as for the scope
    - Skip description for scope for now
    - Allways use SynkzoneSSI as the scope client

Future:
[ ] Let the user open and chose files in a zone for their token
    - Useful for giving access to specific files and folders

[ ] Let the user chose coworkers for the token
    - This could be used to create short lived collaboration between some coworkers and the token holder
