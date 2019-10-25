//
// Copyright: (C) 2019 Nestybox Inc.  All rights reserved.
//

package sysboxFsGrpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	pb "github.com/nestybox/sysbox-ipc/sysboxFsGrpc/protobuf"
	"google.golang.org/grpc"
)

// Container info passed by the client to the server across the grpc channel
type ContainerData struct {
	Id       string
	InitPid  int32
	Hostname string
	Ctime    time.Time
	UidFirst int32
	UidSize  int32
	GidFirst int32
	GidSize  int32
}

//
// Establishes grpc connection to sysbox-fs' remote-end.
//
func connect() *grpc.ClientConn {

	// Set up a connection to the server.
	// TODO: Secure me through TLS.
	conn, err := grpc.Dial(sysboxFsAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to sysbox-fs: %v", err)
		return nil
	}

	return conn
}

func containerDataToPbData(data *ContainerData) (*pb.ContainerData, error) {
	pbTime, err := ptypes.TimestampProto(data.Ctime)
	if err != nil {
		return nil, fmt.Errorf("time conversion error: %v", err)
	}

	return &pb.ContainerData{
		Id:       data.Id,
		InitPid:  data.InitPid,
		Ctime:    pbTime,
		UidFirst: data.UidFirst,
		UidSize:  data.UidSize,
		GidFirst: data.GidFirst,
		GidSize:  data.GidSize,
	}, nil
}

func pbDatatoContainerData(pbdata *pb.ContainerData) (*ContainerData, error) {
	cTime, err := ptypes.Timestamp(pbdata.Ctime)
	if err != nil {
		return nil, fmt.Errorf("time conversion error: %v", err)
	}

	return &ContainerData{
		Id:       pbdata.Id,
		InitPid:  pbdata.InitPid,
		Ctime:    cTime,
		UidFirst: pbdata.UidFirst,
		UidSize:  pbdata.UidSize,
		GidFirst: pbdata.GidFirst,
		GidSize:  pbdata.GidSize,
	}, nil
}

//
// Registers container creation in sysbox-fs. Notice that this
// is a blocking call that can potentially have a minor impact
// on container's boot-up speed.
//
func SendContainerRegistration(data *ContainerData) (err error) {
	var pbData *pb.ContainerData

	// Set up sysbox-fs pipeline.
	conn := connect()
	if conn == nil {
		return fmt.Errorf("failed to connect with sysbox-fs")
	}
	defer conn.Close()

	cntrChanIntf := pb.NewSysboxStateChannelClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	pbData, err = containerDataToPbData(data)
	if err != nil {
		return fmt.Errorf("convertion to protobuf data failed: %v", err)
	}

	_, err = cntrChanIntf.ContainerRegistration(ctx, pbData)
	if err != nil {
		return fmt.Errorf("failed to register container with sysbox-fs: %v", err)
	}

	return nil
}

//
// Unregisters container from sysbox-fs.
//
func SendContainerUnregistration(data *ContainerData) (err error) {
	var pbData *pb.ContainerData

	// Set up sysbox-fs pipeline.
	conn := connect()
	if conn == nil {
		return fmt.Errorf("failed to connect with sysbox-fs")
	}
	defer conn.Close()

	cntrChanIntf := pb.NewSysboxStateChannelClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	pbData, err = containerDataToPbData(data)
	if err != nil {
		return fmt.Errorf("convertion to protobuf data failed: %v", err)
	}

	_, err = cntrChanIntf.ContainerUnregistration(ctx, pbData)
	if err != nil {
		return fmt.Errorf("failed to unregister container with sysbox-fs: %v", err)
	}

	return nil
}

//
// Sends a container-update message to sysbox-fs end. At this point, we are
// only utilizing this message for a particular case, update the container
// creation-time attribute, but this function can serve more general purposes
// in the future.
//
func SendContainerUpdate(data *ContainerData) (err error) {
	var pbData *pb.ContainerData

	// Set up sysbox-fs pipeline.
	conn := connect()
	if conn == nil {
		return fmt.Errorf("failed to connect with sysbox-fs")
	}
	defer conn.Close()

	cntrChanIntf := pb.NewSysboxStateChannelClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	pbData, err = containerDataToPbData(data)
	if err != nil {
		return fmt.Errorf("convertion to protobuf data failed: %v", err)
	}

	_, err = cntrChanIntf.ContainerUpdate(ctx, pbData)
	if err != nil {
		return fmt.Errorf("failed to send container-update message to ",
			"sysbox-fs: %v", err)
	}

	return nil
}