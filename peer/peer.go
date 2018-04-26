package main

import (
	"encoding/gob"
	"fmt"
	"github.com/tcnksm/go-input"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	// "strings"
)

type TSP_header struct {
	Type    byte
	Song_id int
}

type TSP_msg struct {
	Header TSP_header
	Msg    []byte
}

const (
	INIT = iota
	LIST
	PLAY
	PAUSE
)

func init() {
	gob.Register(&TSP_header{})
	gob.Register(&TSP_msg{})
}

func main() {
	args := os.Args[:]
	if len(args) != 3 {
		fmt.Println("Usage: ", args[0], "<port> <filedir>")
		os.Exit(1)
	}

	become_discoverable(args)

	handle_command(args)
	// go start serving songs
	// handle input
}

/**
* send local songs to server,
* tracker other peers now have this ip
 */
func become_discoverable(args []string) {
	songs := get_local_song_info(args[2])

	// Connect to tacker
	tracker, err := net.Dial("tcp", "localhost:"+args[1])
	if err != nil {
		fmt.Println("Error connecting to tracker")
		os.Exit(1)
	}
	defer tracker.Close()

	send_local_song_info(tracker, songs)

	// tracker_port int = args[2]
}

/**
* Returns an string slice
* for all files(songs) in directory "songs"
 */
func get_local_song_info(dir_name string) []string {
	info_files, err := ioutil.ReadDir(dir_name)
	if err != nil {
		fmt.Println("cant read songs")
		os.Exit(1)
	}

	song_info := make([]string, len(info_files))
	for i := 0; i < len(info_files); i++ {
		if path.Ext(info_files[i].Name()) != ".info" {
			continue
		}
		content, _ := ioutil.ReadFile(dir_name + "/" + info_files[i].Name())
		song_info = append(song_info, string(content[:]))
	}
	return song_info
}

/**
* send this peers song list to the
* tracker server using gobs to transfer structs
 */
func send_local_song_info(tracker net.Conn, songs []string) {
	msg_content := ""
	for _, s := range songs {
		msg_content += s
	}

	encoder := gob.NewEncoder(tracker)
	msg_struct := &TSP_msg{TSP_header{Type: 0, Song_id: 0}, []byte(msg_content)}
	encoder.Encode(msg_struct)
}

/**
* handle input command from the user
* LIST - get song list from peers
* PLAY <song id> - play song
* PAUSE - pauses playing of song (buffering continues)
* STOP - stop streaming song
* QUIT - <--
 */
func handle_command(args []string) {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	query := "Select option"
	cmd, err := ui.Select(query, []string{"LIST", "PLAY", "PAUSE", "QUIT"}, &input.Options{
		Loop: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cmd)
	os.Exit(1)

	/*
		// read command

		hdr := new(TSP_header)

		switch strings.ToUpper(cmd) {
		case "LIST":
			// get list from server
			hdr.Type = LIST
			send_list_request(hdr) // TSP_header error
			break
		case "PLAY":
			hdr.Type = PLAY
			break
		case "PAUSE":
			hdr.Type = PAUSE
			break
		case "QUIT":
			// quit
		default:
			fmt.Println("invalid command")
		}
	*/
}

func send_list_request(hdr TSP_header, args []string) { // parameter error
	tracker, err := net.Dial("tcp", "localhost:"+args[1])
	if err != nil {
		fmt.Println("Error connecting to tracker")
		os.Exit(1)
	}
	defer tracker.Close()

	msg_content := ""
	encoder := gob.NewEncoder(tracker)
	msg_struct := &TSP_msg{TSP_header{Type: 0, Song_id: 0}, []byte(msg_content)}
	encoder.Encode(msg_struct)
}

/*
func server_side(peer net.Conn) {

}
*/