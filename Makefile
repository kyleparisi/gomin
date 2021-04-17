run:
	bash run.sh

db:
	./schema/manage.py makemigrations
	./schema/manage.py migrate

test:
	bash test.sh
