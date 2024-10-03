<p align="center"><a href="https://www.codechefvit.com" target="_blank"><img src="https://i.ibb.co/4J9LXxS/cclogo.png" width=160 title="CodeChef-VIT" alt="Codechef-VIT"></a>
</p>

<h2 align="center"> CookOff 9.0 Backend </h2>

<h3 align="center"> Powered By </h2>
<div align="center">
  <table>
    <tr>
      <td>
        <div style="border: 1px solid #ccc; padding: 10px; text-align: center; border-radius: 8px; width: 200px;">
          <a href="https://discord.com/invite/GRc3v6n" target="_blank">
            <img src="https://avatars.githubusercontent.com/u/25365178?s=200&v=4" alt="Judge0" style="width:150px; height:150px; border-radius: 8px;">
            <p>Join Judge0 Discord</p>
          </a>
        </div>
      </td>
      <td>
        <div style="border: 1px solid #ccc; padding: 10px; text-align: center; border-radius: 8px; width: 200px;">
          <a href="https://discord.com/invite/dCq3XhgRXs" target="_blank">
            <img src="https://pbs.twimg.com/profile_images/1742205229104259072/2ISO3o7-_400x400.jpg" alt="Sulu" style="width:150px; height:150px; border-radius: 8px;">
            <p>Join Sulu Discord</p>
          </a>
        </div>
      </td>
    </tr>
  </table>
</div>
<br/>

> CookOff is CodeChef VIT’s flagship competitive coding event that tests the coding skills of developers. This is the backend that powers both the admin and participant portals for CookOff 9.0, serving as the backbone to manage users, questions, test cases, and submissions. Designed for efficiency and scalability, our robust backend simplifies the process of overseeing all competition-related tasks, ensuring smooth operations for both administrators and participants alike. With a focus on user-friendliness, it provides all the necessary tools to facilitate a seamless competitive experience.

---

[![DOCS](https://img.shields.io/badge/Documentation-see%20docs-green?style=flat-square&logo=appveyor)](https://documenter.getpostman.com/view/26244894/2sAXqtbgvt)
[![UI](https://img.shields.io/badge/User%20Interface-Link%20to%20UI-orange?style=flat-square&logo=appveyor)](INSERT_UI_LINK_HERE)

## Features

### User

- Sign Up (`/user/signup`)
- Login (`/user/login`)

### Questions

- Create Question (`/question/create`)
- Get all Questions (`/questions`)
- Get a Question (`/question/{question_id}`)
- Get Question by Round (`/question/round`)
- Update Question (`/question`)
- Delete Question (`/question/{question_id}`)

### Testcases

- Create Testcase (`/testcase`)
- Get Testcases by Question (`/questions/{question_id}/testcases`)
- Get a Testcase (`/testcase/{testcase_id}`)
- Update Testcase (`/testcase/{testcase_id}`)
- Delete Testcase (`/testcase/{testcase_id}`)

### Submissions

- Submit Testcase (`/submit`)
- Run Testcase (`/runcode`)

### Leaderboard

- Get Leaderboard (`/leaderboard`)

## Dependencies

- [atlas](https://atlasgo.io/)
- [docker](https://docs.docker.com/)
- [make](https://www.gnu.org/software/make/manual/make.html)

## Instructions

#### Directions to Install

```sh
$ git clone https://github.com/CodeChefVIT/cookoff-9.0-backend.git
$ cd cookoff-9.0-backend
```

#### Prerequisites

1. Setup atlas
2. Configure env (refer .env.example)
3. Configure Makefile

#### Directions to Run

1. Spin up containers

```sh
$ docker compose up --build -d
```

2. Install sqlc

```sh
$ go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

3. Generate sqlc schema and queries

```sh
$ make generate
```

4. Apply migrations

```sh
$ make apply-schema
```

## Contributors

<table>
	<tr align="center" style="font-weight:bold">
		<td>
		Vedant Matanhelia
		<p align="center">
			<img src = "https://avatars.githubusercontent.com/u/71623796?v=4" width="150" height="150" alt="Vedant Matanhelia">
		</p>
			<p align="center">
				<a href = "https://github.com/Xenomorph07">
					<img src = "https://cdn-icons-png.flaticon.com/512/2111/2111432.png" width="36" height = "36" alt="GitHub"/>
				</a>
			</p>
		</td>
		<td>
		Vaibhav Sijaria
		<p align="center">
			<img src = "https://avatars.githubusercontent.com/u/139199971?v=4" width="150" height="150" alt="Vaibhav Sijaria">
		</p>
			<p align="center">
				<a href = "https://github.com/vaibhavsijaria">
					<img src = "https://cdn-icons-png.flaticon.com/512/2111/2111432.png" width="36" height = "36" alt="GitHub"/>
				</a>
			</p>
		</td>
		<td>
		Aman L
		<p align="center">
			<img src = "https://avatars.githubusercontent.com/u/86644389?v=4" width="150" height="150" alt="Aman L">
		</p>
			<p align="center">
				<a href = "https://github.com/Killerrekt">
					<img src = "https://cdn-icons-png.flaticon.com/512/2111/2111432.png" width="36" height = "36" alt="GitHub"/>
				</a>
			</p>
		</td>
		<td>
		Jothish Kamal
		<p align="center">
			<img src = "https://avatars.githubusercontent.com/u/74227363?v=4" width="150" height="150" alt="Jothish Kamal">
		</p>
			<p align="center">
				<a href = "https://github.com/JothishKamal">
					<img src = "https://cdn-icons-png.flaticon.com/512/2111/2111432.png" width="36" height = "36" alt="GitHub"/>
				</a>
			</p>
		</td>
	</tr>
	<tr>
		<td>
		Soham Mahapatra
		<p align="center">
			<img src = "https://avatars.githubusercontent.com/u/155614230?v=4" width="150" height="150" alt="Soham Mahapatra">
		</p>
			<p align="center">
				<a href = "https://github.com/Soham-Maha">
					<img src = "https://cdn-icons-png.flaticon.com/512/2111/2111432.png" width="36" height = "36" alt="GitHub"/>
				</a>
			</p>
		</td>
		<td>
		Abhinav Anand
		<p align="center">
			<img src = "https://avatars.githubusercontent.com/u/140488187?v=4" width="150" height="150" alt="Abhinav Anand">
		</p>
			<p align="center">
				<a href = "https://github.com/Abhinav-055">
					<img src = "https://cdn-icons-png.flaticon.com/512/2111/2111432.png" width="36" height = "36" alt="GitHub"/>
				</a>
			</p>
		</td>
		<td>
		Aman Singh
		<p align="center">
			<img src = "https://avatars.githubusercontent.com/u/80804989?v=4" width="150" height="150" alt="Abhinav Anand">
		</p>
			<p align="center">
				<a href = "https://github.com/DevloperAmanSingh">
					<img src = "https://cdn-icons-png.flaticon.com/512/2111/2111432.png" width="36" height = "36" alt="GitHub"/>
				</a>
			</p>
		</td>
	</tr>
	
</table>

## License

[![License](http://img.shields.io/:license-mit-blue.svg?style=flat-square)](http://badges.mit-license.org)

<p align="center">
	Made with :heart: by <a href="https://www.codechefvit.com" target="_blank">CodeChef-VIT</a>
</p>
