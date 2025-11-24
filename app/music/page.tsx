import MusicPlayer from "@/components/MusicPlayer"
import SongCard from "@/components/SongCard"
import coverImg from "../../assets/sabrina.png"
import coverImg2 from "../../assets/Charlie XCX.png"
import coverImg3 from "../../assets/vondutch.jpg"
import coverImg4 from "../../assets/Ordinary.jpg"


export default function Music(){
  const singers = ["Sabrina Carpenter", "Charlie XCX", "Alex Warren"]
    return(<>
      <div className="flex flex-col gap-5 items-center justify-center min-h-screen">
          <h1 className="font-mono text-4xl font-bold mb-6 ml-4">
      Now Playing
    </h1>
    <div className="pl-5" ><MusicPlayer ></MusicPlayer></div>
    
    <div className="flex flex-row gap-3 ">
     <SongCard cover={coverImg.src} title="Short n Sweet" artists={singers[0]} releaseDate="Novemeber 2024"></SongCard>
     <SongCard cover={coverImg2.src} title="Von Dutch" artists={singers[1]} releaseDate="Feburary 2025"></SongCard>
   </div>
   <div className="flex flex-row gap-3 ">
     <SongCard cover={coverImg3.src} title="Gone" artists={singers[1]} releaseDate="Feburary 2025"></SongCard>
     <SongCard cover={coverImg4.src} title="Ordinary" artists={singers[1]} releaseDate="March 2025"></SongCard>
   </div>
     </div>
    </>)
}