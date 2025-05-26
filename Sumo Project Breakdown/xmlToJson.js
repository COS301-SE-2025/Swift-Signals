const fs = require('fs');
const xml2js = require('xml2js');

const files = ['network/nodes.nod.xml', 'network/edges.edg.xml'];

files.forEach((file) => {
  const xml = fs.readFileSync(file, 'utf-8');
  xml2js.parseString(xml, (err, result) => {
    if (err) throw err;
    const jsonFile = file.replace('.xml', '.json');
    fs.writeFileSync(jsonFile, JSON.stringify(result, null, 2));
    console.log(`Converted ${file} to ${jsonFile}`);
  });
});
