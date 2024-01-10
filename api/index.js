import express from 'express';
import { middleware as graphbrainz } from 'graphbrainz';

const app = express();

// Use the default options:
app.use('/graphbrainz', graphbrainz());

// // or, pass some options:
// app.use('/graphbrainz', graphbrainz({
// 	client: new MusicBrainz({ ... }),
// 	graphiql: true,
// 	...
// }));

app.listen(3000);