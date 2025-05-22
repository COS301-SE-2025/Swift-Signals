interface IntersectionCardProps {
  id: number;
  name: string;
  location: string;
  lanes: string;
  image?: string;
  onSimulate: (id: number) => void;
  onEdit: (id: number) => void;
  onDelete: (id: number) => void;
}


export default IntersectionCard;
