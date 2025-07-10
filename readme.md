# README

## Project Name: Go-Store

### Design Architecture
This project uses the repository pattern due to minimalistic cyclic import while working, such as DDD or the hexagonal pattern, etc. Also, the repository pattern has been chosen for easy maintenance & faster coding for those faced with the rapid timeline. The use case/service layer will be used to separate the complex logic only (this project does not have heavy business logic).

### Installation
1. Clone the repository:
    ```bash
    git clone https://github.com/farelamo/go-store.git
    ```
2. Navigate to the project directory:
    ```bash
    cd go-store
    ```
3. Build docker
    ```bash
    docker compose up --build -d
    ```
4. Go to the vault and fill the env in path which in your vault secret path
    ```bash
    http://0.0.0:8200
    ```
5. Restart store container
    ```bash
    docker restart store
    ```
6. Go to the documentation (Im using Apidog)
   [Documentation Link](https://app.apidog.com/invite/project?token=s63KVaWBFgMpTNWRM3c73)
   or use this link [Apidog Link](https://go-store.apidog.io)

7. Run to your locally with apidog or export it to your loved api docs

### License
This project is licensed under the [MIT License](LICENSE).

### Contact
For questions or feedback, please reach out to [farelamo4@gmail.com].
