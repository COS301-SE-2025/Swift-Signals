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

const IntersectionCard: React.FC<IntersectionCardProps> = ({
  id,
  name,
  location,
  lanes,
  image,
  onSimulate,
  onEdit,
  onDelete,
}) => {
  
};

export default IntersectionCard;
