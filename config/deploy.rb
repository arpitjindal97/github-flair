set :server_host, ENV["server_host"]
set :server_user, ENV["server_user"]
set :docker_img, ENV["DOCKER_IMG"]
set :server_url, "#{fetch(:server_user)}@#{fetch(:server_host)}"

task :deploy do
	on "#{fetch(:server_url)}" do
		execute "export docker_run_id=$(docker ps -a| grep #{fetch(:docker)}|cut -d ' ' -f 1 ) 
		if [ \"$docker_run_id\" == \"\" ];then exit 0
		else	docker stop $docker_run_id
			docker rm $docker_run_id
			docker rmi $docker_run_id
		fi"
		execute "docker pull #{fetch(:docker_img)} >/dev/null"
		execute "docker run -d -p 80:8080 -v /home/arpit/mongo:/data/db #{fetch(:docker_img)}"

	end
end
