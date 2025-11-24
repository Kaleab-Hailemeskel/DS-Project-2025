import {Card} from "@heroui/card"
import {Image} from "@heroui/image"
import { Button } from "@heroui/button";
import {Slider} from "@heroui/slider"
import {CardBody} from "@heroui/card";

interface Song {
  cover: string;       // URL or path to the album cover image
  title: string;       // Song title
  artists: string;   // One or more artist names
  releaseDate: string; // Release date (ISO string or formatted date)
}

export default function SongCard({ cover, title, artists, releaseDate }:Song) {
  return (
    <Card
      isBlurred
      className="border-none bg-background/70 dark:bg-default-100/50 w-[400px] flex flex-row items-center"
      shadow="sm"
    >
      <CardBody className="flex flex-row items-center gap-6 p-4">
        {/* Album Cover */}
        <Image
          alt="Album cover"
          className="object-cover rounded-md"
          height={120}
          width={120}
          shadow="md"
          src={cover}
        />

        {/* Song Info */}
        <div className="flex flex-col justify-center">
          <h2 className="text-xl font-bold text-foreground/90">{title}</h2>
          <p className="text-sm text-foreground/70">{artists}</p>
          <p className="text-xs text-foreground/50 mt-1">Released: {releaseDate}</p>
        </div>
      </CardBody>
    </Card>
  );
}
