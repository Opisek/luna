use axum::{ Json, Router, routing::get };
use dotenv::dotenv;
use minicaldav::{ self };
use serde::Serialize;
use url::Url;

#[derive(Serialize)]
struct Calendar {
  pub url: Url,
  pub name: String,
  pub color: Option<String>,
}
impl From<&minicaldav::Calendar> for Calendar {
  fn from(calendar: &minicaldav::Calendar) -> Self {
    Calendar {
      url: calendar.url().clone(),
      name: calendar.name().clone(),
      color: calendar.color().cloned(),
    }
  }
}

#[derive(Serialize)]
struct Event {
  name: String,
}

#[derive(Clone)]
struct CaldavSettings {
  url: Url,
  username: String,
  password: String,
}

async fn get_calendars(caldav_settings: CaldavSettings) -> Json<Vec<Calendar>> {
  let agent = ureq::Agent::new();

  let credentials = minicaldav::Credentials::Basic(
    caldav_settings.username.into(),
    caldav_settings.password.into()
  );

  let calendars_raw = minicaldav::get_calendars(agent.clone(), &credentials, &caldav_settings.url).unwrap();
  let calendars: Vec<Calendar> = calendars_raw.iter().map(|c| Calendar::from(c)).collect();

  Json(calendars)
}

async fn get_events() -> Json<Vec<Event>> {
  let events = vec![
    Event { name: "Event2".to_string() },
    Event { name: "Event1".to_string() },
  ];
  Json(events)
}

#[tokio::main]
async fn main() {
  dotenv().ok();
  let caldav_settings = CaldavSettings {
    url: Url::parse(std::env::var("CALDAV_URL").unwrap().as_str()).unwrap(),
    username: std::env::var("CALDAV_USERNAME").unwrap(),
    password: std::env::var("CALDAV_PASSWORD").unwrap(),
  };

  let endpoints = Router::new().route("/calendars", get(move || {get_calendars(caldav_settings)}));

  let app = Router::new().nest("/api", endpoints);
  let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
  axum::serve(listener, app).await.unwrap();
}