# blog_aggregator

You'll need Postgres and Go installed to use this

run go install and it should install gator as a CLI command (if not, it shall install blog_aggregator as a CLI command)

you must create a json file in your home directory named ".gatorconfig.json" and include in it the url for your postgres db as follows:
{"db_url":"postgres url","current_user_name":""}

