CREATE TABLE IF NOT EXISTS artists (
  id serial PRIMARY KEY,
  artist varchar(50)
);



CREATE TABLE IF NOT EXISTS songs (
	id serial PRIMARY KEY,
	artist_id int FOREIGN KEY references artists(id),
	song varchar(100) not null,
	release_date date not null,
	lirycs text not null,
	link varchar(200) not null
);

CREATE INDEX IF NOT EXISTS songs_song_idx ON songs(song);