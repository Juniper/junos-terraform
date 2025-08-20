package netconf

import (
	"bufio"
	"fmt"
	"log"
	"strings"
)

const getGroupStr = `<get-configuration database="committed" format="text" >
  <configuration>
  <groups><name>%s</name></groups>
  </configuration>
</get-configuration>
`

// parseGroupData is a function that cleans up the returned data for generic config groups
func parseGroupData(input string) (reply string, err error) {
	var cfgSlice []string

	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		cfgSlice = append(cfgSlice, scanner.Text())
	}

	var cfgSlice2 []string

	for _, v := range cfgSlice {
		r := strings.NewReplacer("}\\n}\\n", "}", "\\n", "", "/*", "", "*/", "", "</configuration-text>", "")

		d := r.Replace(v)

		// fmt.Println(d)

		cfgSlice2 = append(cfgSlice2, d)
	}

	// Figure out the offset. Each Junos version could give us different stuff, so let's look for the group name
	begin := 0
	end := 0

	for k, v := range cfgSlice2 {
		if v == "groups" {
			begin = k + 4
			break
		}
	}

	// We don't want the very end slice due to config terminations we don't need.
	end = len(cfgSlice2) - 3

	// fmt.Printf("Begin = %v\nEnd = %v\n", begin, end)

	reply = strings.Join(cfgSlice2[begin:end], " ")

	return reply, nil
}

// ReadGroup is a helper function
func (g *GoNCClient) ReadGroup(applygroup string) (string, error) {
	g.Lock.Lock()
	err := g.Driver.Dial()

	if err != nil {
		log.Fatal(err)
	}

	getGroupString := fmt.Sprintf(getGroupStr, applygroup)

	reply, err := g.Driver.SendRaw(getGroupString)
	if err != nil {
		return "", err
	}

	err = g.Driver.Close()

	g.Lock.Unlock()

	if err != nil {
		return "", err
	}

	parsedGroupData, err := parseGroupData(reply.Data)
	if err != nil {
		return "", err
	}

	return parsedGroupData, nil
}

// UpdateRawConfig deletes group data and replaces it (for Update in TF)
func (g *GoNCClient) UpdateRawConfig(applygroup string, netconfcall string, commit bool) (string, error) {

	deleteString := fmt.Sprintf(deleteStr, applygroup, applygroup)

	g.Lock.Lock()
	err := g.Driver.Dial()
	if err != nil {
		log.Fatal(err)
	}

	_, err = g.Driver.SendRaw(deleteString)
	if err != nil {
		errInternal := g.Driver.Close()
		g.Lock.Unlock()
		return "", fmt.Errorf("driver error: %+v, driver close error: %+s", err, errInternal)
	}

	groupString := fmt.Sprintf(groupStrXML, netconfcall)

	reply, err := g.Driver.SendRaw(groupString)
	if err != nil {
		errInternal := g.Driver.Close()
		g.Lock.Unlock()
		return "", fmt.Errorf("driver error: %+v, driver close error: %+s", err, errInternal)
	}

	if commit {
		_, err = g.Driver.SendRaw(commitStr)
		if err != nil {
			errInternal := g.Driver.Close()
			g.Lock.Unlock()
			return "", fmt.Errorf("driver error: %+v, driver close error: %+s", err, errInternal)
		}
	}

	err = g.Driver.Close()

	if err != nil {
		g.Lock.Unlock()
		return "", fmt.Errorf("driver close error: %+s", err)
	}

	g.Lock.Unlock()

	return reply.Data, nil
}

// DeleteConfigNoCommit is a wrapper for driver.SendRaw()
// Does not provide mandatory commit unlike DeleteConfig()
func (g *GoNCClient) DeleteConfigNoCommit(applygroup string) (string, error) {

	deleteString := fmt.Sprintf(deleteStr, applygroup, applygroup)

	g.Lock.Lock()
	err := g.Driver.Dial()
	if err != nil {
		log.Fatal(err)
	}

	reply, err := g.Driver.SendRaw(deleteString)
	if err != nil {
		errInternal := g.Driver.Close()
		g.Lock.Unlock()
		return "", fmt.Errorf("driver error: %+v, driver close error: %+s", err, errInternal)
	}

	output := strings.Replace(reply.Data, "\n", "", -1)

	err = g.Driver.Close()

	if err != nil {
		g.Lock.Unlock()
		return "", fmt.Errorf("driver close error: %+s", err)
	}

	g.Lock.Unlock()

	return output, nil
}

// SendRawConfig is a wrapper for driver.SendRaw()
func (g *GoNCClient) SendRawConfig(netconfcall string, commit bool) (string, error) {

	groupString := fmt.Sprintf(groupStrXML, netconfcall)

	g.Lock.Lock()

	err := g.Driver.Dial()

	if err != nil {
		log.Fatal(err)
	}

	reply, err := g.Driver.SendRaw(groupString)
	if err != nil {
		errInternal := g.Driver.Close()
		g.Lock.Unlock()
		return "", fmt.Errorf("driver error: %+v, driver close error: %+s", err, errInternal)
	}

	if commit {
		_, err = g.Driver.SendRaw(commitStr)
		if err != nil {
			errInternal := g.Driver.Close()
			g.Lock.Unlock()
			return "", fmt.Errorf("driver error: %+v, driver close error: %+s", err, errInternal)
		}
	}

	err = g.Driver.Close()

	if err != nil {
		g.Lock.Unlock()
		return "", err
	}

	g.Lock.Unlock()

	return reply.Data, nil
}

// ReadRawGroup is a helper function
func (g *GoNCClient) ReadRawGroup(applygroup string) (string, error) {
	g.Lock.Lock()
	err := g.Driver.Dial()

	if err != nil {
		log.Fatal(err)
	}

	getGroupXMLString := fmt.Sprintf(getGroupXMLStr, applygroup)

	reply, err := g.Driver.SendRaw(getGroupXMLString)
	if err != nil {
		errInternal := g.Driver.Close()
		g.Lock.Unlock()
		return "", fmt.Errorf("driver error: %+v, driver close error: %+s", err, errInternal)
	}

	err = g.Driver.Close()

	g.Lock.Unlock()

	if err != nil {
		return "", err
	}

	return reply.Data, nil
}
