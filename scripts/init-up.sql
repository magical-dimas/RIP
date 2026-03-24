CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    login VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    is_moderator BOOLEAN DEFAULT FALSE
    );

CREATE TABLE IF NOT EXISTS models (
    model_id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(1000) NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    photo_url VARCHAR(255),
    video VARCHAR(255),
    power FLOAT NOT NULL,
    fuel_usage FLOAT NOT NULL
    );

CREATE TABLE IF NOT EXISTS nuclear_calculations (
    calc_id SERIAL PRIMARY KEY,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    creator_id INTEGER NOT NULL REFERENCES users(user_id),
    forming_date TIMESTAMP,
    finish_date TIMESTAMP,
    moderator_id INTEGER REFERENCES users(user_id),
    description VARCHAR(2000)
    );

CREATE TABLE IF NOT EXISTS modelcalcs (
    calc_id INTEGER NOT NULL REFERENCES nuclear_calculations(calc_id),
    model_id INTEGER NOT NULL REFERENCES models(model_id),
    amount INTEGER NOT NULL DEFAULT 1,
    res_power FLOAT,
    res_fuel FLOAT,
    PRIMARY KEY (calc_id, model_id)
    );

INSERT INTO users (login, password, is_moderator) VALUES
                                                      ('user1', 'pass1', false),
                                                      ('moderator', 'modpass', true)
    ON CONFLICT (login) DO NOTHING;

INSERT INTO models (title, description, is_deleted, photo_url, video, power, fuel_usage) VALUES
                                                                                                                                                             ('РИТМ-200', 'РИТМ-200 - водо-водяной ядерный реактор, предназначенный для установки на ледоколах и перспективных плавучих атомных электростанциях, малых АЭС.', false, 'http://localhost:9000/reactorservice/ritm.jpg', 'ritm.mp4', 55, 45),
                                                                                                                                                             ('HTR-PM', 'HTR-PM - Китайский малый модульный ядерный реактор. Это высокотемпературный газоохлаждаемый реактор четвертого поколения с шаровым топливом, разработанный на основе прототипа HTR-10.', false, 'http://localhost:9000/reactorservice/htr.jpg', 'htr.mp4', 210, 260),
                                                                                                                                                             ('КЛТ-40С', 'КЛТ-40С - Российская плавучая атомная теплоэлектростанция (ПАТЭС) проекта 20870, находящаяся в порту города Певек (Чаунский район, Чукотского автономного округа), самая северная АЭС в мире.', false, 'http://localhost:9000/reactorservice/klt.jpg', 'klt.mp4', 70, 85),
                                                                                                                                                             ('IRIS', 'IRIS - Проект реактора четвертого поколения, разработанный международной командой компаний, лабораторий и университетов при координации компании Westinghouse, призван открыть новые рынки для атомной энергетики и создать мост между технологиями реакторов третьего и четвертого поколений.', false, 'http://localhost:9000/reactorservice/iris.jpg', 'iris.mp4', 335, 320),
                                                                                                                                                             ('ГТ-МГР', 'ГТ-МГР - Российско-американский проект по созданию АЭС на базе высокотемпературного газоохлаждаемого реактора с гелиевым теплоносителем, работающего в прямом газотурбинном цикле.', false, 'http://localhost:9000/reactorservice/gtmgr.jpg', 'gtmgr.mp4', 285, 300);