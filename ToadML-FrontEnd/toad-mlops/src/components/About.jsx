import '../styles/About.css'
import { FaGithub } from 'react-icons/fa';

const teamMembers = [
  { name: 'Horacio Hernandez', image: './toadblue.webp' },
  { name: 'Uri Barajas', image: './toadred.webp' },
  { name: 'Marian Sedano', image: './toadpurple.webp' },
  { name: 'Jose Andres Cota', image: './toadgreen.jpg' },
  { name: 'Roberto Osorno', image: './toadyellow.webp' },
];

export default function About() {
  return (
    <section id="about">
      <h2>About</h2>
      <p id="team-name">Team: ToadML</p>
      <div className="team-grid">
        {teamMembers.map((member, index) => (
          <div className="team-member" key={index}>
            <img src={member.image} alt={member.name} className="member-photo" />
            <p>{member.name}</p>
          </div>
        ))}
      </div>
      <p><strong>Dataset 🗄:</strong> <a href="https://www.kaggle.com/datasets/jonathanoheix/face-expression-recognition-dataset/data" target="_blank">Ver aquí</a></p>
      <p><strong>Github <FaGithub style={{ marginLeft: '0px', marginRight: '4px', verticalAlign: 'middle' }}/>:</strong> <a href="https://github.com/MasterHH123/FaceExpressionClassifierBackend" target="_blank"> Ver aquí</a></p>
    </section>
  );
}
