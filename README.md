# Get Filthy
Go CRUD Application 

# Enviromental Variables
$WEBSPHEREHTML variable must be set to where the directory, in which the HTML template files are kept on your platform. On both mac and linux, I keep them in my go src directory for the project. This seems dumb, and I probably should have structured them more like MVVM or MVC framework, but I had limited time to work with. 

Mac OS X
```
export WEBSPHEREHTML=/Users/youruser/go/src/github.com/bsdpunk/webSphere/
```
Linux
```
export WEBSPHEREHTML=/home/youruser/go/src/github.com/bsdpunk/webSphere/
```

If setting up from scratch you need a user, password and database for mysql. The default is root:ContainerBleed with the database Widgets. Runs on port 8080 by default.



# Run Get Filthy
```
getfilthy
```

## Things I regret
* Not abstracting out the db parts, to prevent repeating myself. 
* Not writing my own ORM. But like also, time you guys, I only had a couple days.
* Not sleeping last night
* Not finding a more elegant solution, than an env variable
* Being a generally shit programmer
* Not adding Delete but come on my dudes, not in your description
